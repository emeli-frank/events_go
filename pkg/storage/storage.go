package storage

import (
	"database/sql"
	"errors"
)

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
