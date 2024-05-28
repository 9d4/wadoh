package http

import (
	"net/http"
	"strconv"

	"github.com/9d4/wadoh/html"
)

func webUsers(s *Server, w http.ResponseWriter, r *http.Request) {
	if err := s.templates.Render(w, r, userFromCtx(r.Context()), "dashboard/users.html", nil); err != nil {
		renderError(w, r, err)
		return
	}
}

func webUsersRows(s *Server, w http.ResponseWriter, r *http.Request) {
	since := 0
	if i, err := strconv.Atoi(r.URL.Query().Get("since")); err == nil {
		since = i
	}
	usrs, err := s.storage.Users.List(20, since)
	if err != nil {
		renderError(w, r, err)
		return
	}

	s.templates.RenderPartial(w, "users/rows.html", html.NewPartialData().
		Set("Users", usrs))
}
