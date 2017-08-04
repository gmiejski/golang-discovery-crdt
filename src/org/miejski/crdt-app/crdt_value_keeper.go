package main

import (
	"org/miejski/domain"
	"sync"
	"org/miejski/crdt"
	"org/miejski/discovery"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"org/miejski/simple_json"
	"time"
)

type CrdtValueKeeper interface {
	Get() crdt.Lwwes
	Reset()
	UpdateChannel() (chan domain.DomainUpdateObject)
	Merge(lwwes crdt.Lwwes)
	synchronize()
}

type crdtValueKeeperImpl struct {
	discovery_client *discovery.DiscoveryClient
	stateKeeper      *domain.DomainKeeper
	update_channel   chan domain.DomainUpdateObject
	lock             sync.Mutex
	client           *http.Client
}

func (c *crdtValueKeeperImpl) synchronize() {
	nodes := (*c.discovery_client).CurrentActiveNodes()
	if len(nodes) == 0 {
		return
	}
	first_node := nodes[0]
	lwwes := c.getState(first_node)
	c.Merge(lwwes)
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
	ccc := (*dd).(*crdt.Lwwes)
	c.updateNotSafe(*ccc)
}

func (c *crdtValueKeeperImpl) updateNotSafe(lwwes crdt.Lwwes) {
	state_keeper := *c.stateKeeper
	state_keeper.Set(lwwes)
}

func CreateSafeValueKeeper(dk domain.DomainKeeper, ds *discovery.DiscoveryClient) CrdtValueKeeper {
	http_client := http.Client{
		Timeout: time.Second * 4,
	}
	channel := make(chan domain.DomainUpdateObject)

	k := crdtValueKeeperImpl{stateKeeper: &dk, update_channel: channel, discovery_client: ds, client: &http_client}
	go func() {
		for {
			x, ok := <-channel
			if !ok {
				break
			}
			keeper := *k.stateKeeper
			k.lock.Lock()
			keeper.Add(x)

			update_object, _ := json.Marshal(toCurrentStateDto(keeper.Get()))

			nodes := (*k.discovery_client).CurrentActiveNodes()

			k.lock.Unlock()
			for i := range nodes {
				k.send(nodes[i], update_object)
			}
		}
	}()
	return &k
}

func (c *crdtValueKeeperImpl) send(node discovery.AppNode, object []byte) {
	fmt.Println("sending lwwes after update to: " + node.Url)

	rq, _ := http.NewRequest("POST", node.Url+"/status/synchronize", bytes.NewBuffer(object))
	rs, err := c.client.Do(rq)
	if err != nil {
		fmt.Println(err)
		//fmt.Println(rs.Status)
	} else {
		fmt.Println(rs.Status)
	}
}

func (c *crdtValueKeeperImpl) getState(node discovery.AppNode) crdt.Lwwes {
	fmt.Println("sending lwwes after update to: " + node.Url)
	rq, _ := http.NewRequest("GET", node.Url+"/status", nil)
	rs, err := c.client.Do(rq)
	if err != nil {
		fmt.Println(err)
	}
	var state CurrentStateDto
	simple_json.Unmarshal(rs.Body, &state)
	lwwes := lwwesFromDto(state)
	return lwwes
}
