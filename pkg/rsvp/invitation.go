package rsvp

import "time"

type Invitation struct {
	ID             int
	Title          string
	Description    string
	IsVirtual      bool
	Address        string
	Link           string
	NumberOfSeats  int
	StartTime      *time.Time
	EndTime        *time.Time
	WelcomeMessage string
	IsPublished    bool
}

type InvitationService interface {
	CreateInvitation(i *Invitation) (int, error)
}
