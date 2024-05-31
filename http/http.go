package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/9d4/wadoh/devices"
	"github.com/9d4/wadoh/html"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog/log"
)

const (
	webLoginPath         = "/login"
	webLogoutPostPath    = "/logout"
	webDevicesPath       = "/devices"
	webDevicesDetailPath = "/devices/{id}"
	webDevicesNewPath    = "/devices/new"
	webDevicesQRPath     = "/devices/qr"
	webDevicesDeletePath = "/devices/{id}"

	webDevicesBlockListPath   = "/devices/list"
	webDevicesBlockDetailPath = "/devices/detail/{id}"

	webDevicesPartialGetStatusPath       = "/devices/{id}/_status"
	webDevicesPartialRenamePath          = "/devices/{id}/rename"
	webDevicesPartialAPIKeyPath          = "/devices/{id}/apikey"
	webDevicesPartialAPIKeyGenPath       = "/devices/{id}/genkey"
	webDevicesPartialSendMessagePath     = "/devices/{id}/send"
	webDevicesPartialSendMessagePostPath = "/devices/{id}/send"

	webUsersPath       = "/users"
	webUsersEditPath   = "/users/{id}"
	webUsersDeletePath = "/users/{id}"
	webUsersRowsPath   = "/users/rows"

	apiDevicesSendMessagePath = "/api/devices/send-message"

	userTokenCookieKey  = "jwt"
	userTokenExpiration = 24 * time.Hour

	userTokenCtxKey ctxKey = "userToken"
	deviceCtxKey    ctxKey = "device"
)

type ctxKey string

type handler func(s *Server, w http.ResponseWriter, r *http.Request)

func Error(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		log.Debug().Caller().Err(err).Send()
		tmpl := &html.ErrorTmpl{
			Code:    "ISE",
			Message: "Internal Server Error",
		}
		tmpl.Render(r.Context(), w)
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
	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	})
}
