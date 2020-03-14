package repository

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

const timeFormat = `"2006-01-02 15:04"`

type NullString struct {
	sql.NullString
}

func NewNullString(value string) NullString {
	return NullString{NullString: sql.NullString{
		String: value,
		Valid:  true,
	}}
}

func (ni NullString) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.String)
}

type NullTime struct {
	sql.NullTime
}

func NewNullTime(value time.Time) NullTime {
	return NullTime{NullTime: sql.NullTime{
		Time:  value,
		Valid: true,
	}}
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid || nt.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(nt.Time.Format(timeFormat)), nil
}

// Time struct wraps time.Time, so can be used with sql and custom marshaling.
type Time struct {
	time.Time
}

func NewTime(value time.Time) Time {
	return Time{Time: value}
}

// Scan implements the Scanner interface.
func (t *Time) Scan(value interface{}) error {
	if value == nil {
		*t = Time{}
		return nil
	}

	*t = Time{Time: value.(time.Time)}
	return nil
}

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

// MarshalJSON marshal time in YYYY-MM-DD HH:MM format.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(t.Format(timeFormat)), nil
}

// After reports whether the time instant t is after u.
func (t Time) After(u Time) bool {
	return t.Time.After(u.Time)
}
