package postgres

import "database/sql"

type Postgres interface {
	DB() *sql.DB
	Tx() (*sql.Tx, error)
	UploadDir() string
}

func New(db *sql.DB, uploadDir string) Postgres {
	return &postgres{db: db, uploadDir: uploadDir}
}

type postgres struct {
	db *sql.DB
	uploadDir string
}

func (p *postgres) DB() *sql.DB {
	return p.db
}

func (p *postgres) Tx() (*sql.Tx, error) {
	return p.db.Begin()
}

func (p *postgres) UploadDir() string {
	return p.uploadDir
}
