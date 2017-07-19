package main

import (
	"org/miejski/domain"
	"sync"
)

type CrdtValueKeeper interface {
	Get() domain.DomainValue
	Reset()
	UpdateChannel() (chan domain.DomainUpdateValue)
}

type crdtValueKeeperImpl struct {
	stateKeeper    *domain.DomainKeeper
	update_channel chan domain.DomainUpdateValue
	lock sync.Mutex
}

func (c *crdtValueKeeperImpl) UpdateChannel() (chan domain.DomainUpdateValue) {
	return c.update_channel
}

func (c *crdtValueKeeperImpl) Reset() {
	state_keeper := *c.stateKeeper
	c.lock.Lock()
	state_keeper.Reset()
	c.lock.Unlock()
}

func (c *crdtValueKeeperImpl) Get() domain.DomainValue {
	c.lock.Lock()
	defer c.lock.Unlock()
	state_keeper := *c.stateKeeper
	return state_keeper.Get()
}

func CreateSafeValueKeeper(dk *domain.DomainKeeper) CrdtValueKeeper {
	channel := make(chan domain.DomainUpdateValue)

	k := crdtValueKeeperImpl{stateKeeper: dk, update_channel: channel}
	go func() {
		for {
			x, ok := <-channel
			if !ok {
				break
			}
			keeper := *k.stateKeeper
			k.lock.Lock()
			keeper.Add(x)
			k.lock.Unlock()
		}
	}()
	return &k
}
