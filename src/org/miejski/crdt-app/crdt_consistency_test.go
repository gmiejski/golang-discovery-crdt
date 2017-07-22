package main

import (
	"testing"
	"org/miejski/domain"
	"sync"
)

var times int = 1000

func TestCrdtConsistency(t *testing.T) {
	// given
	keeper := domain.UnsafeDomainKeeper()
	dk := CreateSafeValueKeeper(&keeper)
	var wg sync.WaitGroup

	// when
	for i := 1; i <= times; i++ {
		wg.Add(1)
		go updateState(&dk, &wg)
	}
	wg.Wait()
	close(dk.UpdateChannel()) // TODO when closing channel it loses single update sometimes? Rework to not enable direct access to channel

	current_value := dk.Get()
	if int(current_value) != times {
		t.Errorf("Inconsistent values! Expected %d but currently %d", times, current_value)
	}

}

func updateState(vk *CrdtValueKeeper, wg *sync.WaitGroup) {
	defer wg.Done()
	(*vk).UpdateChannel() <- domain.DomainUpdateValue(1)
}