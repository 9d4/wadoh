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
	//Redirected indicates that this hx-request has been redirected
	Redirected bool
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

// getHTMX returns *htmxRequest only if that request is not redirected.
// example:
//
//	hxreq ---> /        ---> ask to redirect to /devices
//	hxreq ---> /devices ---> getHTMX(r) == nil (because was redirected from /)
func getHTMX(r *http.Request) *htmxRequest {
	if isHTMXRedirected(r) {
		return nil
	}
	return getHTMXRaw(r)
}

func getHTMXRaw(r *http.Request) *htmxRequest {
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

func setRedirectCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "hx-redirect",
		Value:    "true",
		Path:     "/",
		MaxAge:   1,
		HttpOnly: true,
	})
}

func isHTMXRedirected(r *http.Request) bool {
	c, err := r.Cookie("hx-redirect")
	return err == nil && c.Value == "true"
}

// redirectRefresh redirects and refresh.
func redirectRefresh(w http.ResponseWriter, r *http.Request, path string, code int) {
	hxRequest := getHTMXRaw(r)
	if hxRequest != nil {
		hxRes := htmxResponse{
			HXRedirect: path,
		}
		setHTMXHeaders(w, hxRes)
		setRedirectCookie(w)
		w.WriteHeader(code)
		return
	}

	http.Redirect(w, r, path, code)
}

func redirect(w http.ResponseWriter, r *http.Request, path string, code int) {
	hxRequest := getHTMX(r)
	if hxRequest != nil {
		// here we do full refresh when we cannot determine the baseURL of
		// the origin request.

		hxRes := htmxResponse{
			HXLocation: path,
		}

		u, err := url.Parse(r.Header.Get("Origin"))
		if err == nil {
			hxRes.HXPushUrl = u.JoinPath(path).String()
		} else {
			u, err := url.Parse(r.Referer())
			if err == nil {
				u.RawPath = ""
				hxRes.HXPushUrl = u.JoinPath(path).String()
			} else {
				// or do full refresh instead
				hxRes.HXRefresh = true
			}
		}
		setHTMXHeaders(w, hxRes)
		if !hxRes.HXRefresh {
			setRedirectCookie(w)
		}
		w.WriteHeader(code)
		return
	}

	http.Redirect(w, r, path, code)
}
