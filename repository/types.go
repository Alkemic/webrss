package repository

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (ni NullString) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.String)
}
