package main

import (
	"org/miejski/discovery"
	"org/miejski/domain"
)

func main() {

	discovery_client := discovery.NewDiscoveryClient()

	dk := domain.UnsafeDomainKeeper()

	state_controller := newStateController(&discovery_client, &dk)

	server := NewServer(state_controller)
	server.Start(8080)
}