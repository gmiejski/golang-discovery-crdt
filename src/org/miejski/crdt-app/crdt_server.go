package main

import (
	"net/http"
	"fmt"
	"log"
)

type CrdtServer interface {
	Start(port int)
	Shutdown()
}

type crdtServerImpl struct {
	state_controller StateController
	srv http.Server
}

func NewServer(state_controller StateController) CrdtServer {
	var s http.Server
	server := crdtServerImpl{state_controller, s}
	return &server
}

func (server *crdtServerImpl) Start(port int) {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	http.HandleFunc("/status", GET(server.state_controller.Status))
	http.HandleFunc("/status/update", POST(server.state_controller.Increment))
	http.HandleFunc("/status/reset", POST(server.state_controller.Status))
	fmt.Println("Starting server")
	log.Fatal(srv.ListenAndServe())
}

func (server *crdtServerImpl) Shutdown() {
	server.srv.Shutdown(nil)
}