package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type Jsonb map[string]any

func (j Jsonb) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *Jsonb) Scan(src any) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed.")
	}

	var i any
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*j, ok = i.(map[string]any)
	if !ok {
		return errors.New("type assertion .(map[string]any) failed.")
	}

	return nil
}
