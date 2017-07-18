package domain

type DomainValue int

type DomainKeeper interface {
	Add()
	Get() DomainValue
	Reset()
}

func UnsafeDomainKeeper() DomainKeeper {
	return &unsafeDomainKeeper{}
}

type unsafeDomainKeeper struct {
	value DomainValue
}

func (dk *unsafeDomainKeeper) Add() {
	dk.value += 1
}

func (dk *unsafeDomainKeeper) Get() DomainValue {
	return dk.value
}

func (dk *unsafeDomainKeeper) Reset() {
	dk.value = 0
}