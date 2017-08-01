package main

import (
	"org/miejski/domain"
	"encoding/json"
	"io"
)

type CrdtOperation struct {
	Value     domain.IntElement
	Operation domain.UpdateOperationType
}

func (c *CrdtOperation) Unmarshal(data io.ReadCloser) error {
	decoder := json.NewDecoder(data)
	defer data.Close()
	err := decoder.Decode(c)
	if err != nil {
		return err
	}
	return nil
}

