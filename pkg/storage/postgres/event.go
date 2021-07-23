package postgres

import (
	"database/sql"
	"errors"
	errors2 "events/pkg/errors"
	"events/pkg/events"
	"events/pkg/storage"
	"fmt"
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

func (s *EventStorage) Event(id int) (*events.Event, error) {
	const op = "eventStorage.Event"

	query := fmt.Sprintf("SELECT id, title FROM events WHERE id = %d", id)

	row := s.DB().QueryRow(query)

	var e events.Event
	err := row.Scan(&e.ID, &e.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors2.Wrap(&events.NotFound{Err: err}, op, "scanning into var")
		}
		return nil, errors2.Wrap(err, op, "scanning into var")
	}

	return &e, errors2.Wrap(err, op, "")
}

func (s *EventStorage) Events(uid int) ([]events.Event, error) {
	const op = "eventStorage.Events"

	query := fmt.Sprintf("SELECT id, title FROM events WHERE host_id = %d", uid)

	rows, err := s.DB().Query(query)
	if err != nil {
		return nil, errors2.Wrap(err, op, "executing query")
	}
	defer rows.Close()

	var ee []events.Event
	for rows.Next() {
		var e events.Event
		err = rows.Scan(&e.ID, &e.Title)
		if err != nil {
			return nil, errors2.Wrap(err, op, "scanning into a var")
		}
		ee = append(ee, e)
	}

	return ee, errors2.Wrap(rows.Err(), op, "error after scanning")
}

func (s *EventStorage) SaveEventTx(tx *sql.Tx, title string, uid int) (int, error) {
	const op = "eventStorage.SaveEventTx"

	if tx == nil {
		return 0, errors2.Wrap(storage.TxIsNil, op, "checking transaction")
	}

	query := "INSERT INTO events (title, host_id) VALUES ($1, $2) RETURNING id"
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
		    link = $4,
		    number_of_seats = $5,
		    start_time = $6,
		    end_time = $7,
		    welcome_message = $8,
		    is_published = $9
		WHERE id = $10`

	_, err := tx.Exec(query, i.Title, i.Description, i.IsVirtual, i.Link, i.NumberOfSeats,
		i.StartTime, i.EndTime, i.WelcomeMessage, i.IsPublished, i.ID)
	if err != nil {
		return errors2.Wrap(err, op, "executing query")
	}

	return nil
}

func (s *EventStorage) PublishEvent(id int) error {
	const op = "eventStorage.PublishEvent"

	query := `Update events SET is_published = $1,WHERE id = $2`

	_, err := s.DB().Exec(query, true, id)
	if err != nil {
		return errors2.Wrap(err, op, "executing query")
	}

	return nil
}
