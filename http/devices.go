package http

import "net/http"

func webDevices(s *Server, w http.ResponseWriter, r *http.Request) {
	renderError(w, r, s.templates.Render(w, r, "dashboard/devices.html", nil))
}
