package postgres

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	errors2 "rsvp/pkg/errors"
	"rsvp/pkg/rsvp"
	"rsvp/pkg/storage"
)

func NewInvitationStorage(base Postgres) (*invitationStorage, error) {
	if base == nil {
		return nil, errors.New("base is nil")
	}
	return &invitationStorage{base}, nil
}

type invitationStorage struct {
	Postgres
}

func (s *invitationStorage) SaveInvitationTx(tx *sql.Tx, title string) (int, error) {
	const op = "userStorage.SaveInvitationTx"

	if tx == nil {
		return 0, errors2.Wrap(storage.TxIsNil, op, "checking transaction")
	}

	query := "INSERT INTO invitations (title) VALUES ($1) RETURNING id"
	var id int
	err := tx.QueryRow(query, title).Scan(&id)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			// todo:: handle error due to duplicate
			//if err.Code ==
			return 0, errors2.Wrap(err, op, "executing query")
		}
		return 0, errors2.Wrap(err, op, "executing query")
	}

	return id, nil
}

func (s *invitationStorage) UpdateInvitationTx(tx *sql.Tx, i *rsvp.Invitation) error {
	const op = "userStorage.UpdateInvitationTx"

	if tx == nil {
		return errors2.Wrap(storage.TxIsNil, op, "checking transaction")
	}

	query := `Update invitations 
		SET 
		    title = $1,
		    description = $2,
		    is_virtual = $3,
		    address = $4,
		    link = $5,
		    seat_number = $6,
		    start_time = $7,
		    end_time = $8,
		    welcome_message = $9,
		    is_published = $10
		WHERE id = $11`

	_, err := tx.Exec(query, i.Title, i.Description, i.IsVirtual, i.Address, i.Link, i.NumberOfSeats,
		i.StartTime, i.EndTime, i.WelcomeMessage, i.IsPublished, i.ID)
	if err != nil {
		return errors2.Wrap(err, op, "executing query")
	}

	return nil
}
