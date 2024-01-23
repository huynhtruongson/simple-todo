package field

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type String struct {
	value  string
	status Status
}

func NewString(v string) String {
	return String{
		value:  v,
		status: Present,
	}
}

func NewNullString() String {
	return String{
		status: Null,
	}
}

func (field String) String() string {
	if field.status == Present {
		return field.value
	}
	return ""
}

// implement the database/sql Scan and Value interfaces
func (field *String) Scan(src interface{}) error {
	if src == nil {
		field.status = Null
		return nil
	}
	switch src := src.(type) {
	case []byte:
		field.value = string(src)
		field.status = Present
		return nil
	case string:
		field.value = src
		field.status = Present
		return nil
	}
	return errors.New(fmt.Sprintf("fail to scan data from sql: %s", src))

}

func (field *String) Value() (driver.Value, error) {
	if field.status != Present {
		return nil, nil
	}

	return []byte(field.value), nil
}

func (field *String) UnmarshalJSON(b []byte) error {
	var v *string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v == nil {
		field.status = Null
		return nil
	}
	field.value = *v
	field.status = Present

	return nil
}

func (field String) MarshalJSON() ([]byte, error) {
	if field.status != Present {
		return []byte("null"), nil
	}
	return json.Marshal(field.value)
}
