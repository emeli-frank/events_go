package services

import (
	"database/sql"
	errors2 "rsvp/pkg/errors"
	"rsvp/pkg/rsvp"
	"rsvp/pkg/storage/postgres"
)

type invitationRepo interface {
	postgres.Postgres
	SaveInvitationTx(tx *sql.Tx, title string) (int, error)
	UpdateInvitationTx(tx *sql.Tx, i *rsvp.Invitation) error
}

func NewInvitationService(r invitationRepo) *invitationService {
	return &invitationService{r: r}
}

type invitationService struct {
	r invitationRepo
}

func (s *invitationService) CreateInvitation(i *rsvp.Invitation) (int, error) {
	const op = "userStorage.CreateInvitation"

	tx, err := s.r.Tx()
	if err != nil {
		return 0, errors2.Wrap(err, op, "getting tx")
	}

	id, err := s.r.SaveInvitationTx(tx, i.Title)
	if err != nil {
		_ = tx.Rollback()
		return 0, errors2.Wrap(err, op, "saving invitation title via repo")
	}

	i.ID = id

	err = s.r.UpdateInvitationTx(tx, i)
	if err != nil {
		_ = tx.Rollback()
		return 0, errors2.Wrap(err, op, "updating invitation via repo")
	}

	return id, errors2.Wrap(tx.Commit(), op, "committing")
}
