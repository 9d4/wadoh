package http

import (
	"context"
	"net/http"
	"time"

	"github.com/9d4/wadoh/users"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	webLoginPath      = "/login"
	webDevicesPath    = "/devices"
	webDevicesNewPath = "/devices/new"
	webDevicesQRPath  = "/devices/qr"

	userTokenCookieKey  = "jwt"
	userTokenExpiration = 24 * time.Hour

	userCtxKey      ctxKey = "user"
	userTokenCtxKey ctxKey = "userToken"
)

type ctxKey string

type handler func(s *Server, w http.ResponseWriter, r *http.Request)

func renderError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newCtxUserToken(ctx context.Context, tk jwt.Token) context.Context {
	return context.WithValue(ctx, userTokenCtxKey, tk)
}

func userTokenFromCtx(ctx context.Context) jwt.Token {
	tk, ok := ctx.Value(userTokenCtxKey).(jwt.Token)
	if ok {
		return tk
	}
	return nil
}

func newCtxUser(ctx context.Context, user *users.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func userFromCtx(ctx context.Context) *users.User {
	user, ok := ctx.Value(userCtxKey).(*users.User)
	if ok {
		return user
	}
	return nil
}
