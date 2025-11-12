package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/9d4/wadoh/devices"
	"github.com/9d4/wadoh/internal"
)

const (
	webLoginPath            = "/login"
	webLogoutPostPath       = "/logout"
	webDevicesPath          = "/devices"
	webDevicesAllPath       = "/devices/all"
	webDevicesDetailPath    = "/devices/{id}"
	webDevicesNewPath       = "/devices/new"
	webDevicesQRPath        = "/devices/qr"
	webDevicesDeletePath    = "/devices/{id}"
	webDevicesReconnectPath = "/devices/{id}/reconnect"

	webDevicesBlockDetailPanePath        = "/devices/{id}/pane"
	webDevicesPartialGetStatusPath       = "/devices/{id}/_status"
	webDevicesPartialRenamePath          = "/devices/{id}/rename"
	webDevicesPartialAPIKeyGenPath       = "/devices/{id}/genkey"
	webDevicesPartialSendMessagePostPath = "/devices/{id}/send"
	webDevicesSaveWebhookPostPath        = "/devices/{id}/save_webhook"

	webUsersPath       = "/users"
	webUsersEditPath   = "/users/{id}"
	webUsersDeletePath = "/users/{id}"
	webUsersRowsPath   = "/users/rows"

	apiDevicesSendMessagePath          = "/api/devices/send-message"
	apiDevicesSendMessageImagePath     = "/api/devices/send-message-image"
	apiDevicesSendMessageImageLinkPath = "/api/devices/send-message-image-link"

	userTokenCookieKey  = "jwt"
	userTokenExpiration = 24 * time.Hour

	userTokenCtxKey ctxKey = "userToken"
	deviceCtxKey    ctxKey = "device"
)

type ctxKey string

type handler func(s *Server, w http.ResponseWriter, r *http.Request)

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

func SetFlash(w http.ResponseWriter, value string) {
	maxAge := 0
	if value == "" {
		maxAge = -1
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
	})
}

func SetError(w http.ResponseWriter, err error) {
	maxAge := 0
	value := ""
	if err == nil {
		maxAge = -1
	} else {
		_, value, _ = internal.ParseError(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "error",
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
	})
}
