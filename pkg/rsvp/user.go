package rsvp

type User struct {
	Names string
	Email string
	IsAdmin bool
}

type UserService interface {
	CreateUser(u *User, password string) (int, error)
	EmailMatchPassword(email string, password string) (bool, int, error)
	User(uid int) (*User, error)
}
