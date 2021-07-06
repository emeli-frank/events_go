package postgres

import "database/sql"

type Postgres interface {
	DB() *sql.DB
	Tx() (*sql.Tx, error)
}

func New(db *sql.DB) Postgres {
	return &postgres{db: db}
}

type postgres struct {
	db *sql.DB
}

func (p *postgres) DB() *sql.DB {
	return p.db
}

func (p *postgres) Tx() (*sql.Tx, error) {
	return p.db.Begin()
}
