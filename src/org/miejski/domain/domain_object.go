package domain

import "strconv"

type IntElement struct {
	Value int
}

func (e IntElement) Get() string {
	return strconv.Itoa(e.Value)
}

type UpdateOperationType string

const (
	ADD    UpdateOperationType = "ADD"
	REMOVE UpdateOperationType = "REMOVE"
)

type DomainUpdateObject struct {
	Value     int
	Operation UpdateOperationType
}
