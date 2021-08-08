package postgres

import (
	"database/sql"
	"errors"
	errors2 "events/pkg/errors"
	"events/pkg/events"
	"events/pkg/storage"
	"fmt"
	"github.com/lib/pq"
	"os"
	"path/filepath"
	"strings"
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

func (s *EventStorage) event(db storage.DB, id int) (*events.Event, error) {
	const op = "eventStorage.event"

	if db == nil {
		return nil, errors2.Wrap(errors.New("db is nil"), op, "fetching event")
	}

	query := fmt.Sprintf(
		`SELECT id, title, description, link, start_time, 
					end_time, welcome_message, cover_image_path, is_published, host_id 
				FROM events 
				WHERE id = %d`,
		id,
	)

	row := db.QueryRow(query)

	var description, link, welcomeMessage, coverImagePath sql.NullString
	var startTime, endTime sql.NullTime

	var e events.Event
	err := row.Scan(&e.ID, &e.Title, &description, &link, &startTime,
		&endTime, &welcomeMessage, &coverImagePath, &e.IsPublished, &e.HostID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors2.Wrap(&events.NotFound{Err: err}, op, "scanning into var")
		}
		return nil, errors2.Wrap(err, op, "scanning into var")
	}

	e.Description = storage.NullableStrToStr(description)
	e.Link = storage.NullableStrToStr(link)
	e.WelcomeMessage = storage.NullableStrToStr(welcomeMessage)
	e.CoverImagePath = storage.NullableStrToStr(coverImagePath)
	e.StartTime = storage.SqlTimeToTime(startTime)
	e.EndTime = storage.SqlTimeToTime(endTime)

	return &e, nil
}

func (s *EventStorage) Event(id int) (*events.Event, error) {
	const op = "eventStorage.Event"

	e, err := s.event(s.DB(), id)
	return e, errors2.Wrap(err, op, "getting event")
}

func (s *EventStorage) EventTx(tx *sql.Tx, id int) (*events.Event, error) {
	const op = "eventStorage.EventTx"

	e, err := s.event(tx, id)
	return e, errors2.Wrap(err, op, "getting event")
}

func (s *EventStorage) Events(uid int) ([]events.Event, error) {
	const op = "eventStorage.Events"

	query := fmt.Sprintf(`SELECT id, title, end_time, cover_image_path, left(description, 100) 
						FROM events 
						WHERE host_id = %d 
						ORDER BY created_at DESC`, uid)

	rows, err := s.DB().Query(query)
	if err != nil {
		return nil, errors2.Wrap(err, op, "executing query")
	}
	defer rows.Close()

	var ee []events.Event
	for rows.Next() {
		var e events.Event
		err = rows.Scan(&e.ID, &e.Title, &e.EndTime, &e.CoverImagePath, &e.Description)
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

	query := `UPDATE events 
		SET 
		    title = $1,
		    description = $2,
		    link = $3,
		    start_time = $4,
		    end_time = $5,
		    welcome_message = $6,
		    is_published = $7,
		    cover_image_path = $8
		WHERE id = $9`

	_, err := tx.Exec(query, i.Title, i.Description, i.Link,
		i.StartTime, i.EndTime, i.WelcomeMessage, i.IsPublished, i.CoverImagePath, i.ID)
	if err != nil {
		return errors2.Wrap(err, op, "executing query")
	}

	return nil
}

func (s *EventStorage) DeleteEventTx(tx *sql.Tx, id int) error {
	const op = "eventStorage.DeleteEventTx"

	if tx == nil {
		return errors2.Wrap(storage.TxIsNil, op, "checking transaction")
	}

	query := fmt.Sprintf("DELETE FROM events WHERE id = %d", id)

	_, err := tx.Exec(query)
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

func (s *EventStorage) SaveEventCover(coverImage []byte, key []string, ext string) error {
	const op = "eventStorage.SaveEventCover"

	fullPath := []string{s.UploadDir()}
	var uniquePath []string
	for _, v := range key[:len(key) - 1] {
		uniquePath = append(uniquePath, v)
	}
	fullPath = append(fullPath, uniquePath...)

	// create random directory if not exist
	if _, err := os.Stat(filepath.Join(fullPath...)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Join(fullPath...), os.ModeDir)
		if err != nil {
			return err
		}
	}

	// return error if file exists
	_, err := os.Stat(filepath.Join(filepath.Join(fullPath...), key[2:len(key)][0] + "." + ext))
	if err == nil { // file exists, return error
		return errors2.Wrap(errors.New("file with name already exist"), op, "checking if file exists")
	} else if !os.IsNotExist(err) { // error is not NotExist, return error
		if err != nil {
			return err
		}
	}

	// create empty file
	file2, err := os.OpenFile(
		filepath.Join(filepath.Join(fullPath...),
			key[2:len(key)][0] + "." + ext),
		os.O_WRONLY|os.O_CREATE,
		os.ModeDir,
	)
	if err != nil {
		return err
	}
	defer file2.Close()

	// copy uploaded file byte into newly  created empty file
	_, err = file2.Write(coverImage)
	if err != nil {
		return err
	}

	return nil
}

func (s *EventStorage) DeleteEventCoverPhoto(path string) error {
	const op = "eventStorage.DeleteEventCoverPhoto"

	if strings.TrimSpace(path) == "" {
		return errors.New("path is empty")
	}

	fmt.Println("deleting ", filepath.Join(s.UploadDir(), path))

	return errors2.Wrap(os.Remove(filepath.Join(s.UploadDir(), path)), op, "deleting file")
}

func (s *EventStorage) SaveEventInvitationsTx(tx *sql.Tx, eventID int, emails []string) error {
	const op = "eventStorage.SaveEventInvitationsTX"

	if tx == nil {
		return errors.New("tx is nil")
	}

	if emails == nil {
		return nil
	}

	query := "INSERT INTO event_invitations (event_id, email) VALUES "

	var params []interface{}

	for k, email := range emails {
		query += fmt.Sprintf("(%d, $%d)", eventID, k + 1)
		if len(emails) - k > 1 {
			query += ","
		}
		params = append(params, email)
	}

	_, err := tx.Exec(query, params...)
	return errors2.Wrap(err, op, "executing query")
}

func (s *EventStorage) Invitations(id int, responded bool, accepted bool) ([]events.Invitation, error) {
	const op = "eventStorage.Invitations"

	var query string
	if !responded {
		query = fmt.Sprintf("SELECT event_id, email, has_responded, response, responded_at FROM event_invitations WHERE event_id = %d AND has_responded = false", id)
	} else {
		if accepted {
			query = fmt.Sprintf("SELECT event_id, email, has_responded, response, responded_at FROM event_invitations WHERE event_id = %d AND has_responded = true AND response = true", id)
		} else {
			query = fmt.Sprintf("SELECT event_id, email, has_responded, response, responded_at FROM event_invitations WHERE event_id = %d AND has_responded = true AND response = false", id)
		}
	}

	rows, err := s.DB().Query(query)
	if err != nil {
		return nil, errors2.Wrap(err, op, "executing query")
	}
	defer rows.Close()

	var ii []events.Invitation
	for rows.Next() {
		var i events.Invitation
		err = rows.Scan(&i.EventID, &i.Email, &i.HasResponded, &i.Response, &i.Time)
		if err != nil {
			return nil, errors2.Wrap(err, op, "scanning into variable")
		}
		ii = append(ii, i)
	}

	return ii, errors2.Wrap(rows.Err(), op, "error after scan")
}
