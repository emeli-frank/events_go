package mock

import (
	"database/sql"
	"io"
)

type Mock struct {
	W io.Writer
	DB *sql.DB
}

func (m *Mock) SeedDB() error {
	/*err := storage.ExecScripts(m.DB, m.ScriptPath)
	if err != nil {
		return err
	}*/

	return nil
}
