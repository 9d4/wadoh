package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	redirectContinueParam = "continue"
)

func (s *Server) authenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		redirectLoginPage := func() {
			target := webLoginPath
			if r.URL.Path != "/" {
				target = target + fmt.Sprintf("?%s=%s", redirectContinueParam, r.URL.String())
			}
			http.Redirect(w, r, target, http.StatusFound)
		}

		if err != nil {
			redirectLoginPage()
			return
		}

		if token == nil || jwt.Validate(token, s.tokenAuth.ValidateOptions()...) != nil {
			redirectLoginPage()
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (s *Server) unAuthenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, _, err := jwtauth.FromContext(r.Context())
		if err == nil {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
