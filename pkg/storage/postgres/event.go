package postgres

import (
	"database/sql"
	"errors"
	errors2 "events/pkg/errors"
	"events/pkg/events"
	"events/pkg/storage"
	"github.com/lib/pq"
)

func NewEventStorage(base Postgres) (*EventStorage, error) {
	if base == nil {
		return nil, errors.New("base is nil")
	}
	return &EventStorage{base}, nil
}

type EventStorage struct {
	Postgres
}

func (s *EventStorage) SaveEventTx(tx *sql.Tx, title string, uid int) (int, error) {
	const op = "eventStorage.SaveEventTx"

	if tx == nil {
		return 0, errors2.Wrap(storage.TxIsNil, op, "checking transaction")
	}

	query := "INSERT INTO events (title, host) VALUES ($1, $2) RETURNING id"
	var id int
	err := tx.QueryRow(query, title, uid).Scan(&id)
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

func (s *EventStorage) UpdateEventTx(tx *sql.Tx, i *events.Event) error {
	const op = "eventStorage.UpdateEventTx"

	if tx == nil {
		return errors2.Wrap(storage.TxIsNil, op, "checking transaction")
	}

	query := `Update events 
		SET 
		    title = $1,
		    description = $2,
		    is_virtual = $3,
		    address = $4,
		    link = $5,
		    number_of_seats = $6,
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
