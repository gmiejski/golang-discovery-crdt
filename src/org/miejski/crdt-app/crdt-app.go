package main

import (
	"org/miejski/discovery"
	"net/http"
	"org/miejski/domain"
)

func main() {

	discovery_client := discovery.NewDiscoveryClient()

	dk := domain.UnsafeDomainKeeper()

	state_controller := newStateController(&discovery_client, &dk)

	http.HandleFunc("/status", state_controller.Status)
	http.HandleFunc("/status/update", state_controller.Increment)
	http.ListenAndServe(":8080", nil)
}