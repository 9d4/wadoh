package http

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/9d4/wadoh/wadoh-be/pb"
)

func apiDevicesSendMessage(s *Server, w http.ResponseWriter, r *http.Request) {
	var req *pb.SendMessageRequest
	if err := parseJSON(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Debug().Caller().Err(err).Msg("request parsing err")
		return
	}
    req.Jid = deviceFromCtx(r.Context()).ID

	go func() {
		s.pbCli.SendMessage(context.Background(), req)
	}()
}
