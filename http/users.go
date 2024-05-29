package http

import (
	"net/http"
	"strconv"

	"github.com/9d4/wadoh/html"
)

func webUsers(s *Server, w http.ResponseWriter, r *http.Request) {
	users, err := s.storage.Users.List(20, 0)
	if err != nil {
		renderError(w, r, err)
		return
	}

	tmpl := &html.UsersTmpl{
		Rows: &html.UsersRowsPartial{Users: users},
	}
	renderError(w, r, s.templates.R(r.Context(), w, tmpl))
}

func webUsersRows(s *Server, w http.ResponseWriter, r *http.Request) {
	since := 0
	if i, err := strconv.Atoi(r.URL.Query().Get("since")); err == nil {
		since = i
	}
	users, err := s.storage.Users.List(20, since)
	if err != nil {
		renderError(w, r, err)
		return
	}

	tmpl := &html.UsersRowsPartial{Users: users}

	renderError(w, r, s.templates.R(r.Context(), w, tmpl))
}
