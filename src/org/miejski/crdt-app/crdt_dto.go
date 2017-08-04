package main

import (
	"time"
	"org/miejski/crdt"
	"io"
	"encoding/json"
)

type CurrentStateDto struct {
	AddSet    map[string]time.Time
	RemoveSet map[string]time.Time
}

func (c *CurrentStateDto) Unmarshal(data io.ReadCloser) error {
	decoder := json.NewDecoder(data)
	defer data.Close()
	err := decoder.Decode(c)
	if err != nil {
		return err
	}
	return nil
}

func toCurrentStateDto(lwwes crdt.Lwwes) CurrentStateDto {
	add := map[string]time.Time{}
	remove := map[string]time.Time{}
	times := lwwes.Add_set
	for x, y := range times {
		add[x] = y
	}
	for x, y := range lwwes.Remove_set {
		remove[x] = y
	}
	dto := CurrentStateDto{AddSet: add, RemoveSet: remove}
	return dto
}

func lwwesFromDto(dto CurrentStateDto) crdt.Lwwes {
	add := map[string]time.Time{}
	remove := map[string]time.Time{}
	for x, y := range dto.AddSet {
		add[x] = y
	}
	for x, y := range dto.RemoveSet {
		remove[x] = y
	}
	return crdt.Lwwes{add, remove}
}

type ReadableState struct {
	Values []string
}

