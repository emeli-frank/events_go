package storage

import (
	"database/sql"
	"errors"
	"time"
)

type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var (
	ErrNotAPtr = errors.New("deactivation setter is not a pointer")
	TxIsNil    = errors.New("transaction is nil")
)

const (
	MaxBigIntValue = 18446744073709551615
)

func StrToNullableStr(s1 string) sql.NullString {
	var s2 sql.NullString
	if s1 == "" {
		return s2
	}

	s2.Valid = true
	s2.String = s1
	return s2
}

func NullableStrToStr(s1 sql.NullString) string {
	var s2 string
	if !s1.Valid {
		return s2
	}

	return s1.String
}

func IntToNullableInt(n1 int64) sql.NullInt64 {
	var n2 sql.NullInt64
	if n1 == 0 {
		return n2
	}

	n2.Valid = true
	n2.Int64 = n1
	return n2
}

func NullableIntToInt(n1 sql.NullInt64) int64 {
	var n2 int64
	if !n1.Valid {
		return n2
	}

	return n1.Int64
}

func NullableFloatToFloat(n1 sql.NullFloat64) float64 {
	var n2 float64
	if !n1.Valid {
		return n2
	}

	return n1.Float64
}

// TimeToSQLTime takes *time.Time and returns corresponding sql.NullTime
// that can be safely saved in the DB
func TimeToSQLTime(t1 *time.Time) sql.NullTime {
	var t2 sql.NullTime
	if t1 != nil {
		t2.Valid = true
		t2.Time = *t1
	}

	return t2
}

// SqlTimeToTime takes sq..NullTime and returns corresponding *time.Time
// If time is not zero time, it returns it otherwise it returns a zero time
func SqlTimeToTime(t1 sql.NullTime) *time.Time {
	var t2 *time.Time
	if !t1.Valid {
		return t2
	}

	return &t1.Time
}
