package services

import (
	"database/sql"
	errors2 "events/pkg/errors"
	"events/pkg/events"
	"events/pkg/storage/postgres"
)

type eventRepo interface {
	postgres.Postgres
	SaveEventTx(tx *sql.Tx, title string, uid int) (int, error)
	UpdateEventTx(tx *sql.Tx, i *events.Event) error
}

func NewEventService(r eventRepo) *eventService {
	return &eventService{r: r}
}

type eventService struct {
	r eventRepo
}

func (s *eventService) CreateEvent(e *events.Event, uid int) (int, error) {
	const op = "userStorage.CreateEvent"

	tx, err := s.r.Tx()
	if err != nil {
		return 0, errors2.Wrap(err, op, "getting tx")
	}

	id, err := s.r.SaveEventTx(tx, e.Title, uid)
	if err != nil {
		_ = tx.Rollback()
		return 0, errors2.Wrap(err, op, "saving event title via repo")
	}

	e.ID = id

	err = s.r.UpdateEventTx(tx, e)
	if err != nil {
		_ = tx.Rollback()
		return 0, errors2.Wrap(err, op, "updating event via repo")
	}

	return id, errors2.Wrap(tx.Commit(), op, "committing")
}

/*var TimeRule = func(start, end *time.Time) validation.Rule {
	return validation.RuleFunc(
		func(v interface{}) error {
			now := time.Now()
			if start != nil {
				if start.Before(now) {
					return validation.NewError("start time has passed")
				}
				if end != nil && end.Before(*start) {
					return validation.NewError("end time cannot be set before start time")
				}
			} else {
				if end != nil {
					return validation.NewError("end time cannot be set without a start time")
				}
			}
			return nil
		},
	)
}*/

