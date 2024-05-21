package http

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/9d4/wadoh/users"
)

func webLogin(s *Server, w http.ResponseWriter, r *http.Request) {
	renderError(w, r, s.templates.Render(w, r, "login.html", nil))
}

func webLoginPost(s *Server, w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	user, err := s.storage.Users.GetBy(username)
	if err != nil {
		renderError(w, r, err)
		return
	}
	if err := users.ComparePwd(user.Password, password); err != nil {
		renderError(w, r, err)
		return
	}
	token, tokenString, err := createUserToken(s.tokenAuth, user)
	if err != nil {
		renderError(w, r, err)
		return
	}

	setUserTokenCookie(w, token, tokenString)

	if next := r.FormValue(redirectContinueParam); next != "" {
		url, err := url.Parse(next)
		if err == nil {
			http.Redirect(w, r, url.String(), http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func webLogoutPost(s *Server, w http.ResponseWriter, r *http.Request) {
	clearUserTokenCookie(w)
	header := w.Header()
	header["HX-Refresh"] = []string{"true"}
	w.WriteHeader(http.StatusNoContent)
}

func createUserToken(tokenAuth *jwtauth.JWTAuth, u *users.User) (token jwt.Token, tokenString string, err error) {
	token, tokenString, err = tokenAuth.Encode(map[string]interface{}{
		"id":  u.ID,
		"exp": time.Now().Add(userTokenExpiration).Unix(),
		"sub": u.Username,
	})
	if err != nil {
		return
	}

	return
}

func setUserTokenCookie(w http.ResponseWriter, token jwt.Token, tokenString string) {
	cookie := &http.Cookie{
		Name:     userTokenCookieKey,
		Value:    tokenString,
		Path:     "/",
		Expires:  token.Expiration(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}

func clearUserTokenCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   userTokenCookieKey,
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}