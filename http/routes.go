package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func initializeRoutes(s *Server) {
	handle := func(fn handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fn(s, w, r)
		}
	}

	r := s.router
	r.With(s.unAuthenticated).Get(webLoginPath, handle(webLogin))
	r.With(s.unAuthenticated).Post(webLoginPath, handle(webLoginPost))

	r.Group(func(r chi.Router) {
		r.Use(s.authenticated)

		r.Post(webLogoutPostPath, handle(webLogoutPost))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			redirectRefresh(w, r, webDevicesPath, http.StatusFound)
		})
		r.Get(webDevicesPath, handle(webDevices))
		r.Get(webDevicesDetailPath, handle(webDevicesDetail))
		r.Get(webDevicesBlockDetailPanePath, handle(webDevicesBlockDetailPane))
		r.Get(webDevicesNewPath, handle(webDevicesNew))
		r.Post(webDevicesQRPath, handle(webDevicesQRPost))
		r.Get(webDevicesReconnectPath, handle(webDeviceReconnectView))
		r.Delete(webDevicesDeletePath, handle(webDeviceDelete))

		r.Get(webDevicesPartialGetStatusPath, handle(webDevicesGetStatus))
		r.Get(webDevicesPartialRenamePath, handle(webDevicesRename))
		r.Put(webDevicesPartialRenamePath, handle(webDevicesRenamePut))
		r.Post(webDevicesPartialAPIKeyGenPath, handle(webDevicesPartialAPIKeyGenerate))
		r.Post(webDevicesPartialSendMessagePostPath, handle(webDevicePartialSendMessagePost))
		r.Post(webDevicesSaveWebhookPostPath, handle(webDevicesSaveWebhookPost))
	})

	r.Group(func(r chi.Router) {
		r.Use(s.authenticatedAdmin)

		r.Get(webUsersPath, handle(webUsers))
		r.Post(webUsersPath, handle(webUsersAdd))
		r.Get(webUsersEditPath, handle(webUsersEdit))
		r.Post(webUsersEditPath, handle(webUsersEditPost))
		r.Delete(webUsersDeletePath, handle(webUsersDelete))
	})

	r.Group(func(r chi.Router) {
		r.Use(s.apiAuthenticated)
		r.Post(apiDevicesSendMessagePath, handle(apiDevicesSendMessage))
	})
}
