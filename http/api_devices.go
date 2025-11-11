package http

import (
	"context"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/9d4/wadoh/wadoh-be/pb"
	"github.com/gorilla/schema"
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
