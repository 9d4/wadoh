package http

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/skip2/go-qrcode"

	"github.com/9d4/wadoh/devices"
	"github.com/9d4/wadoh/wadoh-be/pb"
)

type devicesPageData struct {
	Devices []devices.Device
}

func (d devicesPageData) DevicePath(dev devices.Device) string {
	return path.Join(webDevicesPath, dev.ID)
}

func webDevices(s *Server, w http.ResponseWriter, r *http.Request) {
	user := userFromCtx(r.Context())
	devices, err := s.storage.Devices.ListByOwnerID(user.ID)

	if err != nil {
		renderError(w, r, err)
		return
	}
	renderError(w, r, s.templates.Render(w, r, "dashboard/devices.html", devicesPageData{
		Devices: devices,
	}))
}

func webDevicesNew(s *Server, w http.ResponseWriter, r *http.Request) {
	renderError(w, r, s.templates.Render(w, r, "dashboard/devices_new.html", nil))
}

func webDevicesQRPost(s *Server, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	cli, err := s.pbCli.RegisterDevice(r.Context(), &pb.RegisterDeviceRequest{Phone: phone, PushNotification: true})
	if err != nil {
		renderError(w, r, err)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")

	lastQr, lastCode := "", ""
	send := func() {
		html := ""

		if lastCode != "" {
			html += fmt.Sprintf("<code>%s</code>", lastCode)
		}
		if lastQr != "" {
			html += fmt.Sprintf(`<img src="data:image/png;base64, %s"/>`, lastQr)
		}

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
	user := userFromCtx(r.Context())

	s.storage.Devices.Save(&devices.Device{
		ID:       jid,
		Name:     name,
		OwnerID:  user.ID,
		LinkedAt: time.Now(),
	})

	http.Redirect(w, r, webDevicesPath, http.StatusFound)
}

func webDevicesGetStatus(s *Server, w http.ResponseWriter, r *http.Request) {
	jid := chi.RouteContext(r.Context()).URLParam("id")
	user := userFromCtx(r.Context())
	status := "Unknown"

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := s.pbCli.Status(context.Background(), &pb.StatusRequest{
			Jid: jid,
		})
		if err == nil {
			status = statusResponseToString(res)
		}
	}()

	dev, err := s.storage.Devices.GetByID(jid)
	if err != nil {
		// TODO: handle in a better way
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if dev.OwnerID != user.ID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	wg.Wait()
	w.Write([]byte(status))
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
