package main

import (
	"time"
	"org/miejski/crdt"
	"org/miejski/domain"
	"strconv"
)

type CurrentStateDto struct {
	AddSet    map[string]time.Time
	RemoveSet map[string]time.Time
}

func toCurrentStateDto(lwwes crdt.Lwwes) CurrentStateDto {
	add := map[string]time.Time{}
	remove := map[string]time.Time{}
	times := lwwes.Add_set
	for x, y := range times {
		val := x.Get()
		add[val] = y
	}
	for x, y := range lwwes.Remove_set {
		val := x.Get()
		remove[val] = y
	}
	dto := CurrentStateDto{AddSet: add, RemoveSet: remove}
	return dto
}

func lwwesFromDto(dto CurrentStateDto) crdt.Lwwes {
	add := map[crdt.Element]time.Time{}
	remove := map[crdt.Element]time.Time{}
	for x, y := range dto.AddSet {
		v, _ := strconv.Atoi(x)
		add[&domain.IntElement{v}] = y
	}
	for x, y := range dto.RemoveSet {
		v, _ := strconv.Atoi(x)
		remove[&domain.IntElement{v}] = y
	}
	return crdt.Lwwes{add, remove}
}

type ReadableState struct {
	Values []string
}

