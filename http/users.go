package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/9d4/wadoh/html"
	"github.com/9d4/wadoh/users"
)

func webUsers(s *Server, w http.ResponseWriter, r *http.Request) {
	users, err := s.storage.Users.List(20, 0)
	if err != nil {
		Error(w, r, err)
		return
	}

	tmpl := &html.UsersTmpl{
		Rows: &html.UsersRowsBlock{Users: users},
	}
	Error(w, r, s.templates.R(r.Context(), w, tmpl))
}

func webUsersRows(s *Server, w http.ResponseWriter, r *http.Request) {
	since := 0
	if i, err := strconv.Atoi(r.URL.Query().Get("since")); err == nil {
		since = i
	}
	users, err := s.storage.Users.List(20, since)
	if err != nil {
		Error(w, r, err)
		return
	}

	tmpl := &html.UsersRowsBlock{Users: users}

	Error(w, r, s.templates.R(r.Context(), w, tmpl))
}

func webUsersAdd(s *Server, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if err := s.storage.Users.Save(&users.User{
		Name:     name,
		Username: username,
		Password: password,
	}); err != nil {
		http.Redirect(w, r, webUsersPath, http.StatusFound)
		return
	}

	SetFlash(w, fmt.Sprintf("User created with username %s", username))
	http.Redirect(w, r, webUsersPath, http.StatusFound)
}
