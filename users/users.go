package users

import "time"

type User struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Username  string      `json:"username"`
	Password  string      `json:"password"`
	Perm      Permissions `json:"perm"`
	CreatedAt time.Time   `json:"created_at"`
}
