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
}

func (c *crdtValueKeeperImpl) synchronize() {
	nodes := (*c.discovery_client).CurrentActiveNodes()
	if len(nodes) == 0 {
		return
	}
	first_node := nodes[0]
	lwwes := getState(first_node)
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
	channel := make(chan domain.DomainUpdateObject)

	k := crdtValueKeeperImpl{stateKeeper: &dk, update_channel: channel, discovery_client: ds}
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

			fmt.Println("Im in")

			nodes := (*k.discovery_client).CurrentActiveNodes()
			for i := range nodes {
				send(nodes[i], update_object)
			}

			k.lock.Unlock()
		}
	}()
	return &k
}

func send(node discovery.AppNode, object []byte) {
	fmt.Println("sending lwwes after update to: " + node.Url)
	client := http.Client{}
	rq, _ := http.NewRequest("POST", node.Url+"/status/synchronize", bytes.NewBuffer(object))
	rs, err := client.Do(rq)
	if err != nil {
		fmt.Println(err)
		fmt.Println(rs.Status)
	}
	fmt.Println(rs.Status)
}

func getState(node discovery.AppNode) crdt.Lwwes {
	fmt.Println("sending lwwes after update to: " + node.Url)
	client := http.Client{}
	rq, _ := http.NewRequest("GET", node.Url+"/status", nil)
	rs, err := client.Do(rq)
	if err != nil {
		fmt.Println(err)
	}
	var state CurrentStateDto
	simple_json.Unmarshal(rs.Body, &state)
	lwwes := lwwesFromDto(state)
	return lwwes
}