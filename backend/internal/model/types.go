package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StringList 字符串列表，以 JSON 形式存储于 MySQL 的 json 列。
type StringList []string

// Value 实现 driver.Valuer。
func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("marshal string list: %w", err)
	}
	return string(b), nil
}

// Scan 实现 sql.Scanner。
func (s *StringList) Scan(value interface{}) error {
	if value == nil {
		*s = StringList{}
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("unexpected type for string list: %T", value)
	}
	if len(b) == 0 {
		*s = StringList{}
		return nil
	}
	if err := json.Unmarshal(b, s); err != nil {
		return fmt.Errorf("unmarshal string list: %w", err)
	}
	return nil
}
