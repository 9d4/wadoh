package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/9d4/wadoh/html"
	"github.com/9d4/wadoh/users"
	"github.com/go-chi/chi/v5"
)

func webUsers(s *Server, w http.ResponseWriter, r *http.Request) {
	htmx := getHTMX(r)

	since := 0
	if i, err := strconv.Atoi(r.URL.Query().Get("since")); err == nil {
		since = i
	}
	users, err := s.storage.Users.List(20, since)
	if err != nil {
		Error(s, w, r, err)
		return
	}
	if err != nil {
		Error(s, w, r, err)
		return
	}
	rows := &html.UsersRowsBlock{Users: users}

	// if htmx request, send rows only
	if htmx != nil {
		Error(s, w, r, s.templates.R(r.Context(), w, rows))
		return
	}

	tmpl := &html.UsersTmpl{
		List: &html.UsersListBlock{
			Rows: rows,
		},
	}
	Error(s, w, r, s.templates.R(r.Context(), w, tmpl))
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

func webUsersEdit(s *Server, w http.ResponseWriter, r *http.Request) {
	id := chi.RouteContext(r.Context()).URLParam("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		Error(s, w, r, err)
		return
	}

	user, err := s.storage.Users.GetBy(idInt)
	if err != nil {
		Error(s, w, r, err)
		return
	}

	tmpl := &html.UsersTmpl{
		Detail: &html.UsersDetailBlock{
			User: user,
		},
	}
	Error(s, w, r, s.templates.R(r.Context(), w, tmpl))
}

func webUsersEditPost(s *Server, w http.ResponseWriter, r *http.Request) {
	id := chi.RouteContext(r.Context()).URLParam("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		Error(s, w, r, err)
		return
	}

	newUser := &users.User{
		ID:       uint(idInt),
		Name:     r.FormValue("name"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Perm: users.Permissions{
			Admin: r.FormValue("perm.admin") == "on",
		},
	}
	if err := s.storage.Users.Update(newUser); err != nil {
		Error(s, w, r, err)
		return
	}

	SetFlash(w, "User updated succesfully.")
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

func webUsersDelete(s *Server, w http.ResponseWriter, r *http.Request) {
	id := chi.RouteContext(r.Context()).URLParam("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		redirect(w, r, r.Referer(), http.StatusFound)
		return
	}
	_ = idInt
	if err := s.storage.Users.Delete(uint(idInt)); err != nil {
		redirect(w, r, r.Referer(), http.StatusFound)
		return
	}

	SetFlash(w, "User deleted succesfully")
	redirect(w, r, webUsersPath, http.StatusFound)
}
