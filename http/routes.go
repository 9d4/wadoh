package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func initializeRoutes(s *Server) {
	handle := func(fn handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fn(s, w, r)
		}
	}

	r := s.router
	r.With(s.unAuthenticated).Get(webLoginPath, handle(webLogin))
	r.With(s.unAuthenticated).Post(webLoginPath, handle(webLoginPost))

	r.Group(func(r chi.Router) {
		r.Use(s.authenticated)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, webDevicesPath, http.StatusFound)
		})
		r.Get(webDevicesPath, handle(webDevices))
		r.Get(webDevicesNewPath, handle(webDevicesNew))
		r.Post(webDevicesQRPath, handle(webDevicesQRPost))
	})
}
