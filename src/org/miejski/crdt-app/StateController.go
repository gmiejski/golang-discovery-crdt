package main

import (
	"net/http"
	"org/miejski/discovery"
	"org/miejski/domain"
	"fmt"
)

type StateController interface {
	Status(writer http.ResponseWriter, request *http.Request)
}

func newStateController(
	discoveryClient *discovery.DiscoveryClient,
	stateKeeper *domain.DomainKeeper) StateController {

	controller := StateControllerImpl{*discoveryClient, *stateKeeper}
	return &controller
}

type StateControllerImpl struct {
	client      discovery.DiscoveryClient
	stateKeeper domain.DomainKeeper
}

func (c *StateControllerImpl) Status(w http.ResponseWriter, request *http.Request) {
	value := c.stateKeeper.Get()
	fmt.Fprintf(w, "%d", value)
}
