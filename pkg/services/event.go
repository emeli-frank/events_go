package services

import (
	"database/sql"
	errors2 "events/pkg/errors"
	"events/pkg/events"
	"events/pkg/storage/postgres"
	"events/pkg/validation"
	"math/rand"
	"path/filepath"
)

type eventRepo interface {
	postgres.Postgres
	Events(uid int) ([]events.Event, error)
	Event(id int) (*events.Event, error)
	EventTx(tx *sql.Tx, id int) (*events.Event, error)
	SaveEventTx(tx *sql.Tx, title string, uid int) (int, error)
	UpdateEventTx(tx *sql.Tx, i *events.Event) error
	PublishEvent(id int) error
	SaveEventCover(fileBytes []byte, uniquePath []string, ext string) error
}

func NewEventService(r eventRepo) *eventService {
	return &eventService{r: r}
}

type eventService struct {
	r eventRepo
}

func (s *eventService) Events(uid int) ([]events.Event, error) {
	const op = "userStorage.Events"

	ee, err := s.r.Events(uid)

	return ee, errors2.Wrap(err, op, "getting events from repo")
}

func (s *eventService) Event(id int) (*events.Event, error) {
	const op = "userStorage.Event"

	e, err := s.r.Event(id)

	return e, errors2.Wrap(err, op, "getting event from repo")
}

func (s *eventService) CreateEvent(e *events.Event, coverImage []byte, coverImageExt string, uid int) (int, error) {
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

	err = s.updateEventTx(tx, e, coverImage, coverImageExt) // todo::// user service so that images can be updated too
	if err != nil {
		_ = tx.Rollback()
		return 0, errors2.Wrap(err, op, "updating event via repo")
	}

	return id, errors2.Wrap(tx.Commit(), op, "committing")
}

func (s *eventService) UpdateEvent(e *events.Event, coverImage []byte, coverImageExt string) error {
	const op = "eventService.UpdateEvent"

	tx, err := s.r.Tx()
	if err != nil {
		return errors2.Wrap(err, op, "getting tx")
	}

	err = s.updateEventTx(tx, e, coverImage, coverImageExt)
	if err != nil {
		_ = tx.Rollback()
		return errors2.Wrap(err, op, "updating event via update service")
	}

	return errors2.Wrap(tx.Commit(), op, "committing")
}

func (s *eventService) updateEventTx(tx *sql.Tx, e *events.Event, coverImage []byte, coverImageExt string) error {
	const op = "eventService.updateEventTx"

	savedEvent, err := s.r.EventTx(tx, e.ID)
	if err != nil {
		return errors2.Wrap(err, op, "getting event")
	}

	savedEvent.Title = 			e.Title
	savedEvent.Description = 	e.Description
	savedEvent.Link = 			e.Link
	savedEvent.StartTime = 		e.StartTime
	savedEvent.EndTime = 		e.EndTime
	savedEvent.WelcomeMessage = e.WelcomeMessage
	savedEvent.IsPublished = 	e.IsPublished

	// add cover image path if cover image exist
	var key []string
	if coverImage != nil {
		randStrs := RandStringBytes(32)
		key = []string{randStrs[:2], randStrs[2:4], randStrs[4:]}
		savedEvent.CoverImagePath = filepath.Join(key...) + "." + coverImageExt
	}

	err = s.r.UpdateEventTx(tx, savedEvent)
	if err != nil {
		return errors2.Wrap(err, op, "updating event via repo")
	}

	if coverImage != nil {
		err = s.saveEventCover(coverImage, key, coverImageExt)
		if err != nil {
			return errors2.Wrap(err, op, "updating event via repo")
		}
	}

	return nil
}

func (s *eventService) PublishEvent(id int, uid int) error {
	const op = "eventService.PublishEvent"

	e, err := s.Event(id)
	if err != nil {
		return errors2.Wrap(err, op, "getting events")
	}

	// checking if user is authorized to access resource
	if uid != e.HostID {
		return errors2.Wrap(&events.Unauthorized{Err: err}, op,
			"checking if user is authorized to access resource")
	}

	// checking that event has not been published already
	if e.IsPublished {
		return errors2.Wrap(&events.Conflict{Err: err}, op,
			"checking that event has not been published already")
	}

	// validation
	err = validation.MergeErrors(
		validation.Value(e.Title, validation.Required),
		validation.Value(e.Description, validation.Required),
		validation.Value(e.Link, validation.Required),
		validation.Value(e.StartTime, validation.Required),
		validation.Value(e.EndTime, validation.Required),
		validation.Value(e.Description, validation.Required),
	)
	if err != nil {
		return errors2.Wrap(err, op, "validating")
	}

	return errors2.Wrap(s.r.PublishEvent(id), op, "calling repo to publish event")
}

func (s *eventService) saveEventCover(fileBytes []byte, key []string, ext string) error {
	const op = "userStorage.saveEventCover"

	return errors2.Wrap(s.r.SaveEventCover(fileBytes, key, ext), op, "saving file via repo")
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


// todo:: take to somewhere more appropriate
func RandStringBytes(n int) string { // todo:: abstract into somewhere
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

