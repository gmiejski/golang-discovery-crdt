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
	mux := http.NewServeMux()
	mux.HandleFunc("/status", rest.GET(server.state_controller.Status))
	mux.HandleFunc("/status/update", rest.POST(server.state_controller.Increment))
	mux.HandleFunc("/status/reset", rest.POST(server.state_controller.Reset))

	discovery.RegisterDiscoveryEndpoints(&server.discovery_client, mux)
	discovery.CreateDiscoveryHeartbeater(&server.discovery_client).Start(2 * time.Second)
	discovery.NewDeadNodeMarker(&server.discovery_client).StartMarking(5 * time.Second)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}

	fmt.Println(fmt.Sprintf("Starting server at port %d", port))
	log.Fatal(srv.ListenAndServe())
}

func (server *crdtServerImpl) Shutdown() {
	server.srv.Shutdown(nil)
}
