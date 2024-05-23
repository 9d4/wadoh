package http

import (
	"net/http"
	"net/url"
)

func webHTMXRedirect(w http.ResponseWriter, r *http.Request, path string, code int) {
	hxRequest := r.Header.Get("HX-Request") == "true"
	hasReferer := r.Header.Get("Referer") != ""
	hasOrigin := r.Header.Get("Origin") != ""
	if hxRequest && (hasReferer || hasOrigin) {
		baseURL := ""
		u, err := url.Parse(r.Header.Get("Referer"))
		if err == nil {
			baseURL = u.Host
		} else {
			u, err = url.Parse(r.Header.Get("Origin"))
			if err == nil {
				baseURL = u.Host
			}
		}

		if baseURL != "" {
			w.Header().Set("HX-Location", path)
			w.Header().Set("HX-Push-Url", u.JoinPath(path).String())
			w.WriteHeader(code)
			return
		}
	}

	w.Header().Set("Location", path)
	w.Header().Set("HX-Redirect", path)
	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(code)
}

func getHTMXCurrentURL(r *http.Request) string {
	return r.Header.Get("HX-Current-URL")
}
