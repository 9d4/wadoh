package http

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"

	"github.com/9d4/wadoh/wadoh-be/pb"
)

var decoder = schema.NewDecoder()

func apiDevicesSendMessage(s *Server, w http.ResponseWriter, r *http.Request) {
	var req *pb.SendMessageRequest
	if err := parseJSON(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Debug().Caller().Err(err).Msg("request parsing err")
		return
	}
	req.Jid = deviceFromCtx(r.Context()).ID

	go func() {
		_, err := s.pbCli.SendMessage(context.Background(), req)
		if err != nil {
			log.Error().Caller().
				Err(err).
				Str("jid", req.Jid).
				Str("phone", req.Phone).
				Msg("SendMessage rpc error")
			return
		}
	}()
}

func apiDevicesSendMessageImage(s *Server, w http.ResponseWriter, r *http.Request) {
	const MAX_FILE_SIZE = 10 << 20 // 10 MB

	type req struct {
		Phone   string `schema:"phone,required"`
		Caption string `schema:"caption"`
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse multipart form"))
		log.Debug().Caller().Err(err).Msg("request parsing err")
		return
	}

	var body req
	if err := decoder.Decode(&body, r.MultipartForm.Value); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to decode form values"))
		log.Debug().Caller().Err(err).Msg("request decoding err")
		return
	}

	file, fh, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to get image file"))
		log.Debug().Caller().Err(err).Msg("request file retrieval err")
		return
	}

	if fh.Size > MAX_FILE_SIZE {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("file size exceeds limit"))
		log.Debug().Msg("file size exceeds limit")
		return
	}

	image, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to read image file"))
		log.Debug().Caller().Err(err).Msg("request file read err")
		return
	}

	pbReq := &pb.SendImageMessageRequest{
		Jid:     deviceFromCtx(r.Context()).ID,
		Phone:   body.Phone,
		Caption: body.Caption,
		Image:   image,
	}

	go func() {
		_, err := s.pbCli.SendImageMessage(context.Background(), pbReq)
		if err != nil {
			log.Error().Caller().Err(err).Msg("SendImageMessage rpc error")
			return
		}
	}()
}

func apiDevicesSendMessageImageLink(s *Server, w http.ResponseWriter, r *http.Request) {
	type req struct {
		Phone   string `schema:"phone,required"`
		Caption string `schema:"caption"`
		URL     string `schema:"url,required"`
	}

	var body req
	if err := parseJSON(r, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Debug().Caller().Err(err).Msg("request parsing err")
		return
	}

	// Fetch image from url and store to byte slice
	resp, err := http.Get(body.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Debug().Caller().Err(err).Msg("failed to fetch image from url")
		return
	}

	defer resp.Body.Close()
	image, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Debug().Caller().Err(err).Msg("failed to read image from response body")
		return
	}

	pbReq := &pb.SendImageMessageRequest{
		Jid:     deviceFromCtx(r.Context()).ID,
		Phone:   body.Phone,
		Caption: body.Caption,
		Image:   image,
	}

	_, err = s.pbCli.SendImageMessage(context.Background(), pbReq)
	if err != nil {
		log.Error().Caller().Err(err).Msg("SendImageMessage rpc error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to send image message"))
		return
	}
}
