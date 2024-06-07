package http

import (
	"net/http"
	"net/url"
)

// htmxRequest represents htmx request headers from https://htmx.org/docs/#request-headers.
type htmxRequest struct {
	Boosted               bool
	CurrentURL            string
	HistoryRestoreRequest bool
	Prompt                string
	Target                string
	TriggerName           string
	Trigger               string
}

type htmxResponse struct {
	HXLocation           string
	HXPushUrl            string
	HXRedirect           string
	HXRefresh            bool
	HXReplaceUrl         string
	HXReswap             string
	HXRetarget           string
	HXReselect           string
	HXTrigger            string
	HXTriggerAfterSettle string
	HXTriggerAfterSwap   string
}

func setHTMXHeaders(w http.ResponseWriter, h htmxResponse) {
	if h.HXLocation != "" {
		w.Header().Set("HX-Location", h.HXLocation)
	}
	if h.HXPushUrl != "" {
		w.Header().Set("HX-Push-Url", h.HXPushUrl)
	}
	if h.HXRedirect != "" {
		w.Header().Set("HX-Redirect", h.HXRedirect)
	}
	if h.HXRefresh {
		w.Header().Set("HX-Refresh", "true")
	}
	if h.HXReplaceUrl != "" {
		w.Header().Set("HX-Replace-Url", h.HXReplaceUrl)
	}
	if h.HXReswap != "" {
		w.Header().Set("HX-Reswap", h.HXReswap)
	}
	if h.HXRetarget != "" {
		w.Header().Set("HX-Retarget", h.HXRetarget)
	}
	if h.HXReselect != "" {
		w.Header().Set("HX-Reselect", h.HXReselect)
	}
	if h.HXTrigger != "" {
		w.Header().Set("HX-Trigger", h.HXTrigger)
	}
	if h.HXTriggerAfterSettle != "" {
		w.Header().Set("HX-Trigger-After-Settle", h.HXTriggerAfterSettle)
	}
	if h.HXTriggerAfterSwap != "" {
		w.Header().Set("HX-Trigger-After-Swap", h.HXTriggerAfterSwap)
	}
}

func getHTMX(r *http.Request) *htmxRequest {
	if c, err := r.Cookie("redirected"); err == nil && c.Value == "true" {
		return nil
	}
	if r.Header.Get("HX-Request") != "true" {
		return nil
	}
	hx := &htmxRequest{
		Boosted:               r.Header.Get("HX-Boosted") == "true",
		CurrentURL:            r.Header.Get("HX-Current-URL"),
		HistoryRestoreRequest: r.Header.Get("HX-History-Restore-Request") == "true",
		Prompt:                r.Header.Get("HX-Prompt"),
		Target:                r.Header.Get("HX-Target"),
		TriggerName:           r.Header.Get("HX-Trigger-Name"),
		Trigger:               r.Header.Get("HX-Trigger"),
	}
	return hx
}

func redirect(w http.ResponseWriter, r *http.Request, path string, code int) {
	setRedirectCookie := func() {
		http.SetCookie(w, &http.Cookie{
			Name:     "redirected",
			Value:    "true",
			Path:     "/",
			MaxAge:   1,
			HttpOnly: true,
		})
	}

	hxRequest := getHTMX(r)
	hasReferer := r.Referer() != ""
	hasOrigin := r.Header.Get("Origin") != ""
	if hxRequest != nil && (hasReferer || hasOrigin) {
		// baseURL := ""
		u, err := url.Parse(r.Referer())
		if err == nil {
			// baseURL = u.Host
		} else {
			u, err = url.Parse(r.Header.Get("Origin"))
			if err == nil {
				// baseURL = u.Host
			}
		}

		setHTMXHeaders(w, htmxResponse{
			HXLocation: path,
			HXRedirect: path,
			HXPushUrl:  u.JoinPath(path).String(),
			HXRefresh:  true,
		})
		setRedirectCookie()
		w.WriteHeader(code)
		return
	}

	// setHTMXHeaders(w, htmxResponse{
	// 	HXLocation: path,
	// 	HXRedirect: path,
	// 	HXRefresh:  true,
	// })
	setRedirectCookie()
	http.Redirect(w, r, path, code)
}
