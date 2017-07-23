package main

import (
	"net/http"
	"fmt"
	"log"
	"org/miejski/discovery"
	"org/miejski/rest"
	"time"
)

type CrdtServer interface {
	Start(port int)
	Shutdown()
}

type crdtServerImpl struct {
	state_controller StateController
	srv              http.Server
	discovery_client discovery.DiscoveryClient
}

func NewServer(
	state_controller *StateController,
	discovery_client *discovery.DiscoveryClient,
) CrdtServer {

	var s http.Server
	server := crdtServerImpl{*state_controller, s, *discovery_client}
	return &server
}

func (server *crdtServerImpl) Start(port int) {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	http.HandleFunc("/status", rest.GET(server.state_controller.Status))
	http.HandleFunc("/status/update", rest.POST(server.state_controller.Increment))
	http.HandleFunc("/status/reset", rest.POST(server.state_controller.Reset))

	discovery.RegisterDiscoveryEndpoints(&server.discovery_client)
	discovery.CreateDiscoveryHeartbeater(&server.discovery_client).Start(2* time.Second)
	fmt.Println("Starting server")
	log.Fatal(srv.ListenAndServe())
}

func (server *crdtServerImpl) Shutdown() {
	server.srv.Shutdown(nil)
}
