package field

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Int struct {
	value  int
	status Status
}

func NewInt(v int) Int {
	return Int{
		value:  v,
		status: Present,
	}
}

func NewNullInt() Int {
	return Int{
		status: Null,
	}
}

func (field Int) Int() int {
	if field.status == Present {
		return field.value
	}
	return 0
}

// implement the database/sql Scan and Value interfaces
func (field *Int) Scan(src interface{}) error {
	if src == nil {
		field.status = Null
		return nil
	}
	v, ok := src.(int64)
	if !ok {
		return errors.New(fmt.Sprintf("fail to scan data from sql: %s", src))
	}
	field.value = int(v)
	field.status = Present

	return nil
}

func (field *Int) Value() (driver.Value, error) {
	if field.status != Present {
		return nil, nil
	}

	return int64(field.value), nil
}

func (field *Int) UnmarshalJSON(b []byte) error {
	var v *int
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

func (field Int) MarshalJSON() ([]byte, error) {
	if field.status != Present {
		return []byte("null"), nil
	}
	return json.Marshal(field.value)
}
