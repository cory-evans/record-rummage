package models

import (
	"encoding/json"
	"errors"
)

func scanJson(val interface{}, out any) error {
	switch v := val.(type) {
	case []byte:
		return json.Unmarshal(v, out)
	case string:
		return json.Unmarshal([]byte(v), out)
	default:
		return errors.New("invalid type")
	}
}
