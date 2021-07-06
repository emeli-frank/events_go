package rsvp

import "context"

type contextKey string
var ContextKeyUser = contextKey("user")

type User struct {
	Names string
	Email string
	IsAdmin bool
}

type key int

// context stuff
var userKey key

// NewUserContext returns a new Context that carries value u.
func NewUserContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// UserFromContext returns the User value stored in ctx, if any or nil otherwise.
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}

type UserService interface {
	CreateUser(u *User, password string) (int, error)
	EmailMatchPassword(email string, password string) (bool, int, error)
	User(uid int) (*User, error)
}
