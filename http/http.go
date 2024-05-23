package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/9d4/wadoh/devices"
	"github.com/9d4/wadoh/users"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	webLoginPath       = "/login"
	webLogoutPostPath  = "/logout"
	webDevicesPath     = "/devices"
	webDevicesItemPath = "/devices/{id}"
	webDevicesNewPath  = "/devices/new"
	webDevicesQRPath   = "/devices/qr"

	webDevicesPartialListPath      = "/devices/list"
	webDevicesPartialItemPath      = "/devices/item/{id}"
	webDevicesPartialGetStatusPath = "/devices/{id}/_status"
	webDevicesPartialRenamePath    = "/devices/{id}/rename"
	webDevicesPartialAPIKeyPath    = "/devices/{id}/apikey"
	webDevicesPartialAPIKeyGenPath = "/devices/{id}/genkey"

	apiDevicesSendMessagePath = "/api/devices/send-message"

	userTokenCookieKey  = "jwt"
	userTokenExpiration = 24 * time.Hour

	userCtxKey      ctxKey = "user"
	userTokenCtxKey ctxKey = "userToken"
	deviceCtxKey    ctxKey = "device"
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

func newCtxDevice(ctx context.Context, device *devices.Device) context.Context {
	return context.WithValue(ctx, deviceCtxKey, device)
}

func deviceFromCtx(ctx context.Context) *devices.Device {
	dev, ok := ctx.Value(deviceCtxKey).(*devices.Device)
	if ok {
		return dev
	}
	return nil
}

func parseJSON(r *http.Request, to any) error {
	if err := json.NewDecoder(r.Body).Decode(to); err != nil {
		return err
	}
	return nil
}
