package simple_json

import (
	"encoding/json"
	"io"
)

type Unmarshalable interface {
	Unmarshal(data io.ReadCloser) error
}

func Unmarshal(data io.ReadCloser, unmarshalable Unmarshalable) error {
	return unmarshalable.Unmarshal(data)
}

func Marshal(e interface{}) (string, error) {
	val, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(val), nil
}
