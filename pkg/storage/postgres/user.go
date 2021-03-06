package postgres

import (
	"database/sql"
	"errors"
	errors2 "events/pkg/errors"
	"events/pkg/events"
	"fmt"
	"github.com/lib/pq"
)

func NewUserStorage(base Postgres) (*userStorage, error) {
	if base == nil {
		return nil, errors.New("base is nil")
	}
	return &userStorage{base}, nil
}

type userStorage struct {
	Postgres
}

func (s *userStorage) SaveUser(u *events.User, hashedPassword string) (int, error) {
	const op = "userStorage.SaveUser"

	query := "INSERT INTO users (names, email, password) VALUES ($1, $2, $3) RETURNING id"
	var id int
	err := s.DB().QueryRow(query, u.Names, u.Email, hashedPassword).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				return 0, errors2.Wrap(&events.Conflict{Err: pqErr}, op, "executing query")
			}
		}
		return 0, errors2.Wrap(err, op, "executing query")
	}

	return id, nil
}

func (s *userStorage) UserIDAndPasswordByEmail(email string) (int, string, error) {
	const op = "userStorage.UserIDAndPasswordByEmail"

	query := `SELECT id, password FROM users WHERE email = $1`

	var id int
	var password string
	err := s.DB().QueryRow(query, email).Scan(&id, &password)
	if err == sql.ErrNoRows {
		err = &events.NotFound{Err: errors.New("user not found")}
		return 0, op, errors2.Wrap(err, op, "scanning into var")
	} else if err != nil {
		return 0, op, errors2.Wrap(err, op, "scanning into var")
	}

	return id, password, nil
}

func (s userStorage) User(uid int) (*events.User, error) {
	const op = "userStorage.User"

	query := fmt.Sprintf(`SELECT 
				users.id,
				users.names, 
				users.email
			FROM users
			WHERE users.id = %d`, uid)

	var u events.User
	err := s.DB().QueryRow(query).Scan(&u.ID, &u.Names, &u.Email)

	return &u, errors2.Wrap(err, op, "querying rows")
}
