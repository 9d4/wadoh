package http

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/9d4/wadoh/html"
	"github.com/9d4/wadoh/users"
)

func webLogin(s *Server, w http.ResponseWriter, r *http.Request) {
	loginTmpl := &html.LoginTmpl{}
	Error(s, w, r, s.templates.R(r.Context(), w, loginTmpl))
}

func webLoginPost(s *Server, w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	sendErr := func() {
		s.templates.RenderPartial(w, "login/error_msg.html", "Credentials does not match our records.")
	}

	user, err := s.storage.Users.GetBy(username)
	if err != nil {
		sendErr()
		return
	}
	if err := users.ComparePwd(user.Password, password); err != nil {
		sendErr()
		return
	}
	token, tokenString, err := createUserToken(s.tokenAuth, user)
	if err != nil {
		sendErr()
		return
	}

	setUserTokenCookie(w, token, tokenString)

	if hx := getHTMX(r); hx != nil && hx.CurrentURL != "" {
		url, _ := url.Parse(hx.CurrentURL)
		if cont := url.Query().Get(redirectContinueParam); cont != "" {
			redirect(w, r, cont, http.StatusOK)
			return
		}
	}

	redirect(w, r, "/", http.StatusOK)
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
	if token == nil {
		http.SetCookie(w, &http.Cookie{
			Name:   userTokenCookieKey,
			MaxAge: -1,
		})
		return
	}
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
