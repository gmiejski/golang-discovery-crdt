package main

import (
	"net/http"
	"org/miejski/discovery"
	"fmt"
	"org/miejski/domain"
	"time"
	"encoding/json"
	"org/miejski/crdt"
	"org/miejski/simple_json"
)

type StateController interface {
	Status(writer http.ResponseWriter, request *http.Request)
	ReadableStatus(writer http.ResponseWriter, request *http.Request)
	Increment(writer http.ResponseWriter, request *http.Request)
	Reset(writer http.ResponseWriter, request *http.Request)
	SynchronizeData(writer http.ResponseWriter, request *http.Request)
}

type StateControllerImpl struct {
	client      discovery.DiscoveryClient
	stateKeeper CrdtValueKeeper
}

func newStateController(
	discoveryClient *discovery.DiscoveryClient,
	stateKeeper *CrdtValueKeeper) StateController {

	go func() {
		for {
			doEvery(2*time.Second, func(t time.Time) {
				//value := (*stateKeeper).Get()
				//fmt.Println(fmt.Sprintf("Current value : #%v", value))
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

func (c *StateControllerImpl) Status(w http.ResponseWriter, request *http.Request) {
	value := c.stateKeeper.Get()
	converted := toCurrentStateDto(value)
	val, _ := json.Marshal(converted)
	fmt.Fprint(w, string(val))
}

func (c *StateControllerImpl) ReadableStatus(w http.ResponseWriter, request *http.Request) {
	value := c.stateKeeper.Get()
	converted := toReadableState(value)
	json_object, _ := json.Marshal(converted)
	fmt.Fprint(w, string(json_object))
}

func (c *StateControllerImpl) Increment(w http.ResponseWriter, request *http.Request) {
	updateInfo := readUpdateInfo(request)
	update_object := domain.DomainUpdateObject{Value: updateInfo.Value.Value, Operation: updateInfo.Operation}
	c.stateKeeper.UpdateChannel() <- update_object
}

func (c *StateControllerImpl) SynchronizeData(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("Synchronizing data!!!")
	var coming_data CurrentStateDto
	simple_json.Unmarshal(request.Body, &coming_data)
	lwwes := lwwesFromDto(coming_data)
	c.stateKeeper.Merge(lwwes)
}

func toReadableState(lwwes crdt.Lwwes) ReadableState {
	values := make([]string, 0)
	elements := lwwes.Get()
	for _, element := range elements {
		values = append(values, element.Get())
	}
	result := ReadableState{Values: values}
	return result
}

func readUpdateInfo(request *http.Request) CrdtOperation {
	var operation CrdtOperation
	err := simple_json.Unmarshal(request.Body, &operation)
	if err != nil {
		panic(err)
	}
	return operation
}

func (c *StateControllerImpl) Reset(w http.ResponseWriter, request *http.Request) {
	c.stateKeeper.Reset()
}
