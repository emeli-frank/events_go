package mock

import (
	"database/sql"
	"events/pkg/storage/postgres"
)

type psql struct{}

func (f *psql) DB() *sql.DB {
	return nil
}

func (f *psql) Tx() (*sql.Tx, error) {
	return nil, nil
}

var Psql postgres.Postgres = &psql{}
