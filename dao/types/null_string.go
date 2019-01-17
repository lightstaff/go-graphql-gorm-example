package types

import (
	"database/sql"
	"encoding/json"
)

// sql.NullStringのラッパー
type NullString struct {
	sql.NullString
}

// MarshalJSON
func (s NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON
func (s *NullString) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	s.String = str
	s.Valid = str != ""
	return nil
}

// 新規作成
func NewNullString(value string) NullString {
	return NullString{
		NullString: sql.NullString{
			String: value,
			Valid:  value != "",
		},
	}
}
