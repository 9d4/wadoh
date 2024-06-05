package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/9d4/wadoh/devices"
	"github.com/9d4/wadoh/html"
	"github.com/9d4/wadoh/internal"
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

	webDevicesBlockListPath       = "/devices/list"
	webDevicesBlockDetailPath     = "/devices/detail/{id}"
	webDevicesBlockDetailPanePath = "/devices/detail/{id}/pane"

	webDevicesPartialGetStatusPath       = "/devices/{id}/_status"
	webDevicesPartialRenamePath          = "/devices/{id}/rename"
	webDevicesPartialAPIKeyGenPath       = "/devices/{id}/genkey"
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

func Error(s *Server, w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	log.Debug().Caller().Err(err).Send()

	// parse http related error
	tmpl := &html.ErrorTmpl{}
	e := parseError(err)
	if e != nil {
		tmpl.Code = e.Code()
		tmpl.Message = e.Error()
		tmpl.Status = errorKindToStatus(e.Kind())
	} else {
		kind, message, code := internal.ParseError(err)
		tmpl.Code = code
		tmpl.Message = message
		tmpl.Status = errorKindToStatus(kind)
	}

	w.WriteHeader(tmpl.Status)
	err = s.templates.R(r.Context(), w, tmpl)
	if err != nil {
		fmt.Fprintf(w, "Something went wrong and unable to render the page. Here's some messages: %s. Code: %s", tmpl.Message, tmpl.Code)
		return
	}
}

var errorKindStatus = map[internal.ErrorKind]int{
	internal.EINTERNAL: http.StatusInternalServerError,
	internal.ENOTFOUND: http.StatusNotFound,
    internal.EBADINPUT: http.StatusBadRequest,
}

func errorKindToStatus(kind internal.ErrorKind) int {
	code, ok := errorKindStatus[kind]
	if ok {
		return code
	}
	return http.StatusInternalServerError
}

// parseError parses error in http layer
func parseError(err error) (e *internal.Error) {
	var strconvErr *strconv.NumError
	if errors.As(err, &strconvErr) {
		return errBadRequest
	}

	return
}

var (
	errBadRequest = internal.NewError(
		internal.EBADINPUT,
		"Your input seems incorrect, please check before try again",
		"bad_request",
	)
)

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
