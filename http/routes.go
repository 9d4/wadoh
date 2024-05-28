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
			w.Header().Set("HX-Redirect", webDevicesPath)
			http.Redirect(w, r, webDevicesPath, http.StatusFound)
		})
		r.Get(webDevicesPath, handle(webDevices))
		r.Get(webDevicesItemPath, handle(webDevices))
		r.Get(webDevicesPartialListPath, handle(webDevicesPartialList))
		r.Get(webDevicesPartialItemPath, handle(webDevicesPartialItem))
		r.Get(webDevicesNewPath, handle(webDevicesNew))
		r.Post(webDevicesQRPath, handle(webDevicesQRPost))
		r.Delete(webDevicesDeletePath, handle(webDeviceDelete))

		r.Get(webDevicesPartialGetStatusPath, handle(webDevicesGetStatus))
		r.Get(webDevicesPartialRenamePath, handle(webDevicesRename))
		r.Put(webDevicesPartialRenamePath, handle(webDevicesRenamePut))
		r.Get(webDevicesPartialAPIKeyPath, handle(webDevicesPartialAPIKey))
		r.Post(webDevicesPartialAPIKeyGenPath, handle(webDevicesPartialAPIKeyGenerate))

		r.Post(webDevicesPartialSendMessagePostPath, handle(webDevicePartialSendMessagePost))
		r.Get(webDevicesPartialSendMessagePath, handle(webDevicePartialSendMessage))
	})

	r.Group(func(r chi.Router) {
		r.Use(s.authenticatedAdmin)

		r.Get(webUsersPath, handle(webUsers))
		r.Get(webUsersRowsPath, handle(webUsersRows))
	})

	r.Group(func(r chi.Router) {
		r.Use(s.apiAuthenticated)
		r.Post(apiDevicesSendMessagePath, handle(apiDevicesSendMessage))
	})
}
