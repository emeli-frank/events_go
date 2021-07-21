package mock

import (
	"database/sql"
	"events/pkg/events"
)

func NewEventStorage() *eventStorage {
	return &eventStorage{
		psql: psql{},
	}
}

type eventStorage struct {
	psql

	DBFn func() *sql.DB
	DBInvoked bool

	TxFn func() (*sql.Tx, error)
	TxInvoked bool

	SaveEventTxFn      func(tx *sql.Tx, title string, uid int) (int, error)
	SaveEventTxInvoked bool

	UpdateEventTxFn func(tx *sql.Tx, i *events.Event) error
	UpdateEventTxInvoked bool
}

func (s *eventStorage) DB() *sql.DB {
	s.DBInvoked = true
	return s.DBFn()
}

func (s *eventStorage) Tx() (*sql.Tx, error) {
	s.TxInvoked = true
	return s.TxFn()
}

func (s *eventStorage) SaveEventTx (tx *sql.Tx, title string, uid int) (int, error) {
	s.SaveEventTxInvoked = true
	return s.SaveEventTxFn(tx, title, uid)
}

func (s *eventStorage) UpdateEventTx (tx *sql.Tx, i *events.Event) error {
	s.UpdateEventTxInvoked = true
	return s.UpdateEventTxFn(tx, i)
}
