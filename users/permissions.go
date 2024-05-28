package users

import "time"

type Permissions struct {
	Admin     bool      `json:"admin"`
	UpdatedAt time.Time `json:"updated_at"`
}
