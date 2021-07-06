package services

import (
	"golang.org/x/crypto/bcrypt"
	errors2 "rsvp/pkg/errors"
	"rsvp/pkg/rsvp"
)

type repo interface {
	SaveUser(u *rsvp.User, hashedPassword string) (int, error)
	UserIDAndPasswordByEmail(email string) (int, string, error)
	User(uid int) (*rsvp.User, error)
}

func NewUserService(r repo) *userService {
	return &userService{r: r}
}

type userService struct {
	r repo
}

func (s *userService) CreateUser(u *rsvp.User, password string) (int, error) {
	const op = "userStorage.CreateUser"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	id, err := s.r.SaveUser(u, string(hash))

	return id, errors2.Wrap(err, op, "saving user via repo")
}

func (s *userService) EmailMatchPassword(email string, password string) (bool, int, error) {
	op := "userService.EmailMatchPassword"

	// todo:: validate email and password

	uid, hashedPassword, err := s.r.UserIDAndPasswordByEmail(email)
	if err != nil {
		switch errors2.Unwrap(err).(type) {
		case *rsvp.NotFound:
			return false, 0, nil
		default:
			return false, 0, err
		}
	}

	// compare user provided and stored password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, 0, nil
	} else if err != nil {
		return false, 0, errors2.Wrap(err, op, "hashing password")
	}

	return true, uid, nil
}

func (s *userService) User(uid int) (*rsvp.User, error) {
	const op = "userService.User"

	u, err := s.r.User(uid)

	return u, errors2.Wrap(err, op, "getting user from repo")
}
