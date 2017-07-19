package main

import (
	"org/miejski/domain"
)

type CrdtValueKeeper interface {
	Get() domain.DomainValue
	Reset()
	UpdateChannel() (chan domain.DomainUpdateValue)
}

type crdtValueKeeperImpl struct {
	stateKeeper    *domain.DomainKeeper
	update_channel chan domain.DomainUpdateValue
}

func (c *crdtValueKeeperImpl) UpdateChannel() (chan domain.DomainUpdateValue) {
	return c.update_channel
}

func (c *crdtValueKeeperImpl) Reset() {
	state_keeper := *c.stateKeeper
	state_keeper.Reset()
}

func (c *crdtValueKeeperImpl) Get() domain.DomainValue {
	state_keeper := *c.stateKeeper
	return state_keeper.Get()
}

func CreateSafeValueKeeper(dk *domain.DomainKeeper) CrdtValueKeeper {
	channel := make(chan domain.DomainUpdateValue)

	k := crdtValueKeeperImpl{dk, channel}
	go func() {
		for {
			x := <-channel
			keeper := *k.stateKeeper
			keeper.Add(x)
		}
	}()
	return &k
}
