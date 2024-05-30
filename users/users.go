package users

import (
	"context"
	"time"
)

type User struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Username  string      `json:"username"`
	Password  string      `json:"password"`
	Perm      Permissions `json:"perm"`
	CreatedAt time.Time   `json:"created_at"`
}

type ctxKey string

const (
	userCtxKey ctxKey = "user"
)

func NewUserContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func UserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userCtxKey).(*User)
	if ok {
		return user
	}
	return nil
}
