package main

import (
	"net/http"
	"org/miejski/discovery"
	"fmt"
	"org/miejski/domain"
	"time"
	"encoding/json"
)

type StateController interface {
	Status(writer http.ResponseWriter, request *http.Request)
	Increment(writer http.ResponseWriter, request *http.Request)
	Reset(writer http.ResponseWriter, request *http.Request)
}

func newStateController(
	discoveryClient *discovery.DiscoveryClient,
	stateKeeper *CrdtValueKeeper) StateController {

	go func() {
		for {
			doEvery(2*time.Second, func(t time.Time) {
				value := (*stateKeeper).Get()
				fmt.Println(fmt.Sprintf("Current value : #%v", value))
			})
		}
	}()

	controller := StateControllerImpl{*discoveryClient, *stateKeeper}
	return &controller
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

type StateControllerImpl struct {
	client      discovery.DiscoveryClient
	stateKeeper CrdtValueKeeper
}

func (c *StateControllerImpl) Status(w http.ResponseWriter, request *http.Request) {
	value := c.stateKeeper.Get()
	converted := toCurrentStateDto(value)
	fmt.Println(converted)
	val, err := json.Marshal(converted)

	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(val))
}

func (c *StateControllerImpl) Increment(w http.ResponseWriter, request *http.Request) {
	updateInfo := readUpdateInfo(request)
	update_object := domain.DomainUpdateObject{Value: updateInfo.Value, Operation: updateInfo.Operation}
	c.stateKeeper.UpdateChannel() <- update_object
}
func readUpdateInfo(request *http.Request) CrdtOperation {
	decoder := json.NewDecoder(request.Body)
	var operation CrdtOperation
	err := decoder.Decode(&operation)
	if err != nil {
		panic(err)
	}
	defer request.Body.Close()
	return operation
}

func (c *StateControllerImpl) Reset(w http.ResponseWriter, request *http.Request) {
	c.stateKeeper.Reset()
}

type CrdtOperation struct {
	Value     domain.IntElement
	Operation domain.UpdateOperationType
}