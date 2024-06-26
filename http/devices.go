package http

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"

	"github.com/9d4/wadoh/devices"
	"github.com/9d4/wadoh/html"
	"github.com/9d4/wadoh/users"
	"github.com/9d4/wadoh/wadoh-be/pb"
)

func webDevices(s *Server, w http.ResponseWriter, r *http.Request) {
	htmx := getHTMX(r)

	user := users.UserFromContext(r.Context())
	devices, err := s.storage.Devices.ListByOwnerID(user.ID)
	if err != nil {
		Error(s, w, r, err)
		return
	}

	listBlock := &html.DevicesListBlock{
		Devices: devices,
	}

	if htmx != nil {
		Error(s, w, r, s.templates.R(r.Context(), w, listBlock))
		return
	}

	tmpl := &html.DevicesTmpl{List: listBlock}
	Error(s, w, r, s.templates.R(r.Context(), w, tmpl))
}

func webDevicesDetail(s *Server, w http.ResponseWriter, r *http.Request) {
	htmx := getHTMX(r)
	id := chi.RouteContext(r.Context()).URLParam("id")
	dev, err := getDevice(s, r.Context(), id)
	if err != nil {
		Error(s, w, r, err)
		return
	}

	detailBlock := &html.DevicesDetailBlock{
		Device:     dev,
		DetailPane: newDevicesDetailPaneBlock(r, dev),
	}

	if htmx != nil {
		Error(s, w, r, s.templates.R(r.Context(), w, detailBlock))
		return
	}

	tmpl := &html.DevicesTmpl{Detail: detailBlock}
	Error(s, w, r, s.templates.R(r.Context(), w, tmpl))
}

func newDevicesDetailPaneBlock(
	r *http.Request,
	dev *devices.Device,
) *html.DevicesDetailPaneBlock {
	queryTab := r.URL.Query().Get("tab")
	return &html.DevicesDetailPaneBlock{
		Device:        dev,
		SubAPIKey:     queryTab == "api",
		SubTryMessage: queryTab == "try_message",
		SubWebhook:    queryTab == "webhook",
		SubMore:       queryTab == "more",
	}
}

func webDevicesBlockDetailPane(s *Server, w http.ResponseWriter, r *http.Request) {
	id := chi.RouteContext(r.Context()).URLParam("id")
	dev, err := getDevice(s, r.Context(), id)
	if err != nil {
		Error(s, w, r, err)
		return
	}

	tmpl := newDevicesDetailPaneBlock(r, dev)
	Error(s, w, r, s.templates.R(r.Context(), w, tmpl))
}

func webDevicesNew(s *Server, w http.ResponseWriter, r *http.Request) {
	tmpl := &html.DevicesNewTmpl{}
	Error(s, w, r, s.templates.R(r.Context(), w, tmpl))
}

func webDevicesQRPost(s *Server, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	cli, err := s.pbCli.RegisterDevice(r.Context(), &pb.RegisterDeviceRequest{Phone: phone, PushNotification: true})
	if err != nil {
		Error(s, w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, r)
		return
	}

	lastQr, lastCode := "", ""
	send := func() {
		html := ""

		if lastCode != "" {
			html += fmt.Sprintf("<code>%s</code>", lastCode)
		}
		if lastQr != "" {
			html += fmt.Sprintf(`<img src="data:image/png;base64, %s"/>`, lastQr)
		}
		html += "\n"

		w.Write([]byte(html))
		flusher.Flush()
	}

	sendSuccess := func() {
		// by sending "success" string, the client will understand and
		// redirect to /devices
		w.Write([]byte("success"))
		flusher.Flush()
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case <-cli.Context().Done():
			return
		default:
		}

		res, err := cli.Recv()
		if err != nil {
			return
		}

		if res.Qr != nil {
			b, err := qrcode.Encode(*res.Qr, qrcode.Medium, 320)
			if err == nil {
				lastQr = base64.RawStdEncoding.EncodeToString(b)
			}
		}

		if res.PairCode != nil {
			lastCode = *res.PairCode
		}

		if res.LoggedIn != nil && res.Jid != nil {
			webDevicesConnectedHandle(s, w, r, name, *res.Jid)
			sendSuccess()
			return
		}

		send()
	}
}

func webDevicesConnectedHandle(s *Server, w http.ResponseWriter, r *http.Request, name, jid string) {
	user := users.UserFromContext(r.Context())

	if err := s.storage.Devices.Save(&devices.Device{
		ID:       jid,
		Name:     name,
		OwnerID:  user.ID,
		LinkedAt: time.Now(),
	}); err != nil {
		log.Debug().Caller().Err(err).Send()
	}

	http.Redirect(w, r, webDevicesPath, http.StatusFound)
}

func webDevicesGetStatus(s *Server, w http.ResponseWriter, r *http.Request) {
	jid := chi.RouteContext(r.Context()).URLParam("id")
	user := users.UserFromContext(r.Context())
	statusString := "Unknown"
	status := pb.StatusResponse_STATUS_UNKNOWN

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := s.pbCli.Status(context.Background(), &pb.StatusRequest{
			Jid: jid,
		})
		if err == nil {
			statusString = statusResponseToString(res)
			status = res.Status
		}
	}()

	dev, err := s.storage.Devices.GetByID(jid)
	if err != nil {
		// TODO: handle in a better way
		w.WriteHeader(http.StatusInternalServerError)
		log.Debug().Err(err).Caller().Send()
		return
	}
	if dev.OwnerID != user.ID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	wg.Wait()
	s.templates.RenderPartial(w, "devices/status.html", map[string]interface{}{
		"Status":       status,
		"StatusString": statusString,
	})
}

func statusResponseToString(res *pb.StatusResponse) string {
	switch res.Status {
	case pb.StatusResponse_STATUS_ACTIVE:
		return "Active"
	case pb.StatusResponse_STATUS_DISCONNECTED:
		return "Disconnected"
	case pb.StatusResponse_STATUS_NOT_FOUND:
		return "Not Found"
	case pb.StatusResponse_STATUS_UNKNOWN:
		return "Unknown"
	default:
		return ""
	}
}

func webDevicesRename(s *Server, w http.ResponseWriter, r *http.Request) {
	jid := chi.RouteContext(r.Context()).URLParam("id")
	user := users.UserFromContext(r.Context())

	device, err := s.storage.Devices.GetByID(jid)
	if err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	if device.OwnerID != user.ID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	s.templates.RenderPartial(w, "devices/rename.html", map[string]interface{}{
		"ID":   device.ID,
		"Name": device.Name,
	})
}

func webDevicesRenamePut(s *Server, w http.ResponseWriter, r *http.Request) {
	ctx := chi.RouteContext(r.Context())
	jid := ctx.URLParam("id")
	user := users.UserFromContext(r.Context())
	newName := r.FormValue("new_name")

	device, err := s.storage.Devices.GetByID(jid)
	if err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	if device.OwnerID != user.ID {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	s.storage.Devices.Rename(device.ID, newName)

	device, err = s.storage.Devices.GetByID(jid)
	if err != nil {
		http.Error(w, "error: Please refresh", http.StatusOK)
		return
	}

	s.templates.RenderPartial(w, "devices/name.html", map[string]interface{}{
		"ID":   device.ID,
		"Name": device.Name,
	})
}

func webDevicesPartialAPIKeyGenerate(s *Server, w http.ResponseWriter, r *http.Request) {
	device, err := getDevice(s, r.Context(), chi.RouteContext(r.Context()).URLParam("id"))
	if err != nil {
		http.Error(w, "unable to render this part", http.StatusOK)
		log.Debug().Caller().Err(err).Send()
		return
	}
	if err := s.storage.Devices.GenNewDevAPIKey(device.ID); err != nil {
		log.Debug().Caller().Err(err).Send()
	}
	device, _ = getDevice(s, r.Context(), chi.RouteContext(r.Context()).URLParam("id"))
	s.templates.RenderPartial(w, "devices/api_key.html", device)
}

func webDevicePartialSendMessagePost(s *Server, w http.ResponseWriter, r *http.Request) {
	phone := r.FormValue("phone")
	message := r.FormValue("message")
	device, err := getDevice(s, r.Context(), chi.RouteContext(r.Context()).URLParam("id"))
	if err != nil {
		http.Error(w, "Permission Denied", http.StatusOK)
		log.Debug().Caller().Err(err).Send()
		return
	}

	go func() {
		s.pbCli.SendMessage(context.Background(), &pb.SendMessageRequest{
			Jid:   device.ID,
			Phone: phone,
			Body:  message,
		})
	}()
	w.Write([]byte("OK"))
}

func webDeviceDelete(s *Server, w http.ResponseWriter, r *http.Request) {
	device, err := getDevice(s, r.Context(), chi.RouteContext(r.Context()).URLParam("id"))
	if err != nil {
		http.Error(w, "Permission Denied", http.StatusOK)
		log.Debug().Caller().Err(err).Send()
		return
	}
	if err := s.storage.Devices.Delete(device.ID); err != nil {
		redirect(w, r, webDevicesPath, http.StatusFound)
		log.Debug().Caller().Err(err).Send()
		return
	}
	redirect(w, r, webDevicesPath, http.StatusFound)
}

func webDevicesSaveWebhookPost(s *Server, w http.ResponseWriter, r *http.Request) {
	device, err := getDevice(s, r.Context(), chi.RouteContext(r.Context()).URLParam("id"))
	if err != nil {
		Error(s, w, r, err)
		return
	}
	err = s.storage.Devices.SaveWebhook(&devices.DeviceWebhook{
		DeviceID: device.ID,
		URL:      r.FormValue("url"),
	})
	if err != nil {
		SetError(w, err)
	} else {
		SetFlash(w, "Webhook url saved successfully")
	}

	redirect(w, r, r.Referer(), http.StatusFound)
}

func getDevice(s *Server, ctx context.Context, deviceID string) (*devices.Device, error) {
	user := users.UserFromContext(ctx)
	device, err := s.storage.Devices.GetByID(deviceID)
	if err != nil {
		return nil, err
	}
	if device.OwnerID != user.ID {
		return nil, os.ErrPermission
	}
	return device, nil
}
