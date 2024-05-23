package http

import "net/http"

func (s *Server) apiAuthenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
            return
		}

		dev, err := s.storage.Devices.GetByAPIToken(token)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
            return
		}
		r = r.WithContext(newCtxDevice(r.Context(), dev))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
