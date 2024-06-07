package http

import (
	"errors"
	"net/http"

	"github.com/9d4/wadoh/html"
	"github.com/9d4/wadoh/users"
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
			//if r.URL.Path != "/" {
			//	target = target + fmt.Sprintf("?%s=%s", redirectContinueParam, r.URL.String())
			//}

			// this prevents from infinite redirect
			setUserTokenCookie(w, nil, "")
			redirect(w, r, target, http.StatusFound)
		}

		if err != nil {
			redirectLoginPage()
			return
		}

		if token == nil || jwt.Validate(token, s.tokenAuth.ValidateOptions()...) != nil {
			redirectLoginPage()
			return
		}

		user, err := s.storage.Users.GetBy(token.Subject())
		if err != nil {
			redirectLoginPage()
			return
		}

		r = r.WithContext(newCtxUserToken(r.Context(), token))
		r = r.WithContext(users.NewUserContext(r.Context(), user))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (s *Server) unAuthenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, _, err := jwtauth.FromContext(r.Context())
		if err == nil {
			redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (s *Server) authenticatedAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := users.UserFromContext(r.Context())
		if !user.Perm.Admin {
			http.Error(w, "You don't have permission to access the resource.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
	return s.authenticated(http.HandlerFunc(fn))
}

func loadFlash(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie("flash"); err == nil && c.Value != "" {
			r = r.WithContext(html.NewFlashContext(r.Context(), c.Value))
			SetFlash(w, "")
		}

		// load error. error is also flash
		if c, err := r.Cookie("error"); err == nil && c.Value != "" {
			r = r.WithContext(html.NewErrorContext(r.Context(), errors.New(c.Value)))
			SetError(w, nil)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
