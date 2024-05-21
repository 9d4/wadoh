package http

import (
	"net/http"
	"net/url"
)

func webHTMXRedirect(w http.ResponseWriter, r *http.Request, path string, code int) {
	if origin := r.Header.Get("Origin"); origin != "" {
		u, err := url.Parse(origin)
		if err == nil {
			w.Header().Set("HX-Location", path)
			w.Header().Set("HX-Push-Url", u.JoinPath(path).String())
			w.WriteHeader(code)
			return
		}
	}

	w.Header().Set("HX-Redirect", path)
	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(code)
}

func getHTMXCurrentURL(r *http.Request) string {
	return r.Header.Get("HX-Current-URL")
}
