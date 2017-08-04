package main

import (
	"testing"
	"org/miejski/domain"
	"sync"
	"github.com/stretchr/testify/assert"
	"org/miejski/discovery"
)

var times int = 1000

func TestCrdtConsistency(t *testing.T) {
	// given
	discovery_client := discovery.NewDiscoveryClient("host", "")
	keeper := domain.UnsafeDomainKeeper()
	dk := CreateSafeValueKeeper(keeper, &discovery_client)
	var wg sync.WaitGroup

	// when
	for i := 1; i <= times; i++ {
		wg.Add(1)
		go updateState(&dk, &wg)
	}
	wg.Wait()
	close(dk.UpdateChannel()) // TODO when closing channel it loses single update sometimes? Rework to not enable direct access to channel

	current_value := dk.Get()
	assert.True(t, current_value.Contains("1"))
}

func updateState(vk *CrdtValueKeeper, wg *sync.WaitGroup) {
	defer wg.Done()
	(*vk).UpdateChannel() <- domain.DomainUpdateObject{1, domain.ADD}
}