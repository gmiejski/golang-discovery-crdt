package domain

import (
	"time"
	"org/miejski/crdt"
)

type DomainKeeper interface {
	Add(DomainUpdateObject)
	Get() crdt.Lwwes
	Reset()
}

func UnsafeDomainKeeper() DomainKeeper {
	return &unsafeDomainKeeper{crdt.CreateLwwes()}
}

type unsafeDomainKeeper struct {
	value crdt.Lwwes
}

func (dk *unsafeDomainKeeper) Add(val DomainUpdateObject) {
	switch val.Operation {
	case ADD:
		{
			value := (*dk).value
			value.Add(&val.Value, time.Now())
		}
	case REMOVE:
		{
			value := (*dk).value
			value.Remove(&val.Value, time.Now())
		}
	}
}

func (dk *unsafeDomainKeeper) Get() crdt.Lwwes {
	return dk.value
}

func (dk *unsafeDomainKeeper) Reset() {
	dk.value = crdt.CreateLwwes()
}
