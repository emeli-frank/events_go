package mock

import (
	"database/sql"
	"io"
	"path/filepath"
	"events/pkg/storage"
)

type Mock struct {
	W io.Writer
	DB *sql.DB
	ScriptPath string
}

func (m *Mock) SeedDB() error {
	err := storage.ExecScripts(m.DB, filepath.Join(m.ScriptPath, "data.sql"))
	if err != nil {
		return err
	}

	return nil
}
