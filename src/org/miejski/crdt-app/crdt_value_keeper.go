package main

import (
	"org/miejski/domain"
	"sync"
	"org/miejski/crdt"
)

type CrdtValueKeeper interface {
	Get() crdt.Lwwes
	Reset()
	UpdateChannel() (chan domain.DomainUpdateObject)
	Merge(lwwes crdt.Lwwes)
}

type crdtValueKeeperImpl struct {
	stateKeeper    *domain.DomainKeeper
	update_channel chan domain.DomainUpdateObject
	lock sync.Mutex
}

func (c *crdtValueKeeperImpl) UpdateChannel() (chan domain.DomainUpdateObject) {
	return c.update_channel
}

func (c *crdtValueKeeperImpl) Reset() {
	state_keeper := *c.stateKeeper
	c.lock.Lock()
	state_keeper.Reset()
	c.lock.Unlock()
}

func (c *crdtValueKeeperImpl) Get() crdt.Lwwes {
	c.lock.Lock()
	defer c.lock.Unlock()
	state_keeper := *c.stateKeeper
	return state_keeper.Get()
}

func (c *crdtValueKeeperImpl) Merge(new crdt.Lwwes) {
	c.lock.Lock()
	defer c.lock.Unlock()

	state_keeper := *c.stateKeeper
	old := state_keeper.Get()
	merge := (&old).Merge(&new)
	d := crdt.LastWriteWinsElementSet(merge)
	dd := &d
	ccc:= (*dd).(*crdt.Lwwes)
	c.updateNotSafe(*ccc)
}

func (c *crdtValueKeeperImpl) updateNotSafe(lwwes crdt.Lwwes) {
	state_keeper := *c.stateKeeper
	state_keeper.Set(lwwes)
}

func CreateSafeValueKeeper(dk domain.DomainKeeper) CrdtValueKeeper {
	channel := make(chan domain.DomainUpdateObject)

	k := crdtValueKeeperImpl{stateKeeper: &dk, update_channel: channel}
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
