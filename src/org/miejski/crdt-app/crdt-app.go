package main

import (
	"org/miejski/discovery"
	"org/miejski/domain"
	"flag"
	"fmt"
)


func main() {

	var port = flag.Int("port", 8080, "Port to bind this server at")
	var joinAddress = flag.String("join", "", "Single node address of the cluster to join")

	flag.Parse()

	this_server_url := fmt.Sprintf("http://localhost:%d", *port)
	discovery_client := discovery.NewDiscoveryClient(this_server_url, *joinAddress)
	keeper := domain.UnsafeDomainKeeper()
	dk := CreateSafeValueKeeper(keeper, &discovery_client)
	dk.synchronize()
	state_controller := newStateController(&discovery_client, &dk)

	server := NewServer(&state_controller, &discovery_client)
	server.Start(*port)
}
