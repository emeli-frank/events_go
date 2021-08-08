package events

import "time"

type Invitation struct {
	EventID int
	Email string
	HasResponded bool
	Response bool
	Time *time.Time
}
