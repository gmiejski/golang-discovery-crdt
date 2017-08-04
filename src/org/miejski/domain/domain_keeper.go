package domain

import (
	"time"
	"org/miejski/crdt"
	"strconv"
)

type DomainKeeper interface {
	Add(DomainUpdateObject)
	Get() crdt.Lwwes
	Reset()
	Set(lwwes crdt.Lwwes)
}

type unsafeDomainKeeper struct {
	value crdt.Lwwes
}

func (dk *unsafeDomainKeeper) Set(lwwes crdt.Lwwes) {
	dk.value = lwwes
}

func UnsafeDomainKeeper() DomainKeeper {
	return &unsafeDomainKeeper{crdt.CreateLwwes()}
}

func (dk *unsafeDomainKeeper) Add(val DomainUpdateObject) {
	switch val.Operation {
	case ADD:
		{
			value := (*dk).value
			value.Add(strconv.Itoa(val.Value), time.Now())
		}
	case REMOVE:
		{
			value := (*dk).value
			value.Remove(strconv.Itoa(val.Value), time.Now())
		}
	}
}

func (dk *unsafeDomainKeeper) Get() crdt.Lwwes {
	return dk.value
}

func (dk *unsafeDomainKeeper) Reset() {
	dk.value = crdt.CreateLwwes()
}
