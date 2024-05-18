package http

import (
	"net/http"
	"time"
)

const (
	webLoginPath = "/login"

	userTokenCookieKey  = "jwt"
	userTokenExpiration = 24 * time.Hour
)

type handler func(s *Server, w http.ResponseWriter, r *http.Request)

func renderError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
