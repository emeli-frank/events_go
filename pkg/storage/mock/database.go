package mock

import (
	"database/sql"
	"events/pkg/storage"
)

// NewDB returns db for test database
func NewDB() (*sql.DB, error) {
	return storage.OpenDB(
		"postgres",
		"host=localhost port=5432 user=events password=password dbname=events_test sslmode=disable",
	)
}
