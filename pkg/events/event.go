package events

import "time"

type Event struct {
	ID             int
	Title          string
	Description    string
	IsVirtual      bool
	Link           string
	NumberOfSeats  int
	StartTime      *time.Time
	EndTime        *time.Time
	WelcomeMessage string
	IsPublished    bool
	HostID         int
}

type EventService interface {
	Events(uid int) ([]Event, error)
	Event(id int) (*Event, error)
	CreateEvent(i *Event, uid int) (int, error)
	PublishEvent(id, uid int) error
}
