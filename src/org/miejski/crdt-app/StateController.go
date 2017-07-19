package main

import (
	"net/http"
	"org/miejski/discovery"
	"fmt"
	"org/miejski/domain"
)

type StateController interface {
	Status(writer http.ResponseWriter, request *http.Request)
	Increment(writer http.ResponseWriter, request *http.Request)
	Reset(writer http.ResponseWriter, request *http.Request)
}

func methodCheck(method string, handler func(w http.ResponseWriter, request *http.Request)) func(w http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		if method == request.Method {
			handler(w,request)
		} else {
			error_msg := fmt.Sprintf("Method not supported: %s, required: %s", request.Method, method)
			http.Error(w, error_msg, http.StatusMethodNotAllowed)
		}
	}
}

func POST(fn func(w http.ResponseWriter, request *http.Request)) func(w http.ResponseWriter, request *http.Request) {
	return methodCheck("POST", fn )
}

func GET(fn func(w http.ResponseWriter, request *http.Request)) func(w http.ResponseWriter, request *http.Request) {
	return methodCheck("GET", fn )
}

func newStateController(
	discoveryClient *discovery.DiscoveryClient,
	stateKeeper *CrdtValueKeeper) StateController {

	controller := StateControllerImpl{*discoveryClient, *stateKeeper}
	return &controller
}

type StateControllerImpl struct {
	client      discovery.DiscoveryClient
	stateKeeper CrdtValueKeeper
}

func (c *StateControllerImpl) Status(w http.ResponseWriter, request *http.Request) {
	value := c.stateKeeper.Get()
	fmt.Fprintf(w, "%d", value)
}

func (c *StateControllerImpl) Increment(w http.ResponseWriter, request *http.Request) {
	c.stateKeeper.UpdateChannel() <- domain.DomainUpdateValue(1)
}

func (c *StateControllerImpl) Reset(w http.ResponseWriter, request *http.Request) {
	c.stateKeeper.Reset()
}
