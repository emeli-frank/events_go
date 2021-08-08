package events

import (
	"time"
)

type Event struct {
	ID             int
	Title          string
	Description    string
	Link           string
	StartTime      *time.Time
	EndTime        *time.Time
	WelcomeMessage string
	CoverImagePath string
	IsPublished    bool
	HostID         int
}

type EventService interface {
	Events(uid int) ([]Event, error)
	Event(id int) (*Event, error)
	CreateEvent(i *Event, invitations []string, coverImage []byte, coverImageExt string, uid int) (int, error)
	UpdateEvent(e *Event, invitations []string, coverImage []byte, coverImageExt string) error
	DeleteEvent(id int) error
	Invitations(id int, responded bool, accepted bool) ([]Invitation, error)
	PublishEvent(id, uid int) error
}
