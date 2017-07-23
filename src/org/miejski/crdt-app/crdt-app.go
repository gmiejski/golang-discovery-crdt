package main

import (
	"org/miejski/discovery"
	"org/miejski/domain"
	"flag"
	"fmt"
)

var port = flag.Int("port", 8080, "Port to bind this server at")

func main() {
	flag.Parse()

	this_server_url := fmt.Sprintf("http://localhost:%d", *port)
	discovery_client := discovery.NewDiscoveryClient(this_server_url)
	keeper := domain.UnsafeDomainKeeper()
	dk := CreateSafeValueKeeper(&keeper)

	state_controller := newStateController(&discovery_client, &dk)

	server := NewServer(&state_controller, &discovery_client)
	server.Start(*port)
}