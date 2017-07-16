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

	http.HandleFunc("/status", GET(state_controller.Status))
	http.HandleFunc("/status/update", POST(state_controller.Increment))
	//http.HandleFunc("/status/reset", POST(state_controller.Status))
	http.ListenAndServe(":8080", nil)
}