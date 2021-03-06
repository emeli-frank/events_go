package services

import (
	"golang.org/x/crypto/bcrypt"
	errors2 "events/pkg/errors"
	"events/pkg/events"
	"events/pkg/storage/postgres"
)

type userRepo interface {
	postgres.Postgres
	SaveUser(u *events.User, hashedPassword string) (int, error)
	UserIDAndPasswordByEmail(email string) (int, string, error)
	User(uid int) (*events.User, error)
}

func NewUserService(r userRepo) *userService {
	return &userService{r: r}
}

type userService struct {
	r userRepo
}

func (s *userService) CreateUser(u *events.User, password string) (int, error) {
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
		case *events.NotFound:
			return false, 0, nil
		default:
			return false, 0, err // todo:: wrap errors
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

func (s *userService) User(uid int) (*events.User, error) {
	const op = "userService.User"

	u, err := s.r.User(uid)

	return u, errors2.Wrap(err, op, "getting user from repo")
}
