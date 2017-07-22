package main

import (
	"org/miejski/discovery"
	"org/miejski/domain"
)

func main() {

	discovery_client := discovery.NewDiscoveryClient()
	keeper := domain.UnsafeDomainKeeper()
	dk := CreateSafeValueKeeper(&keeper)

	state_controller := newStateController(&discovery_client, &dk)

	server := NewServer(state_controller, discovery_client)
	server.Start(8080)
}