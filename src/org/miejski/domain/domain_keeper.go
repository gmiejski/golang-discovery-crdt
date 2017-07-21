package domain

type DomainValue int
type DomainUpdateValue int

type DomainKeeper interface {
	Add(DomainUpdateValue)
	Get() DomainValue
	Reset()
}

func UnsafeDomainKeeper() DomainKeeper {
	return &unsafeDomainKeeper{}
}

type unsafeDomainKeeper struct {
	value DomainValue
}

func (dk *unsafeDomainKeeper) Add(val DomainUpdateValue) {
	dk.value += DomainValue(val)
}

func (dk *unsafeDomainKeeper) Get() DomainValue {
	return dk.value
}

func (dk *unsafeDomainKeeper) Reset() {
	dk.value = 0
}