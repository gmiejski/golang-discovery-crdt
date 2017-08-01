package main

import (
	"net/http"
	"org/miejski/discovery"
	"fmt"
	"org/miejski/domain"
	"time"
	"encoding/json"
	"org/miejski/crdt"
	"strconv"
	"org/miejski/simple_json"
)

type StateController interface {
	Status(writer http.ResponseWriter, request *http.Request)
	ReadableStatus(writer http.ResponseWriter, request *http.Request)
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
	val, _ := json.Marshal(converted)
	fmt.Fprint(w, string(val))
}

func (c *StateControllerImpl) ReadableStatus(w http.ResponseWriter, request *http.Request) {
	value := c.stateKeeper.Get()
	converted := toReadableState(value)
	json_object, _ := json.Marshal(converted)
	fmt.Fprint(w, json_object)
}

func (c *StateControllerImpl) Increment(w http.ResponseWriter, request *http.Request) {
	updateInfo := readUpdateInfo(request)
	update_object := domain.DomainUpdateObject{Value: updateInfo.Value.Value, Operation: updateInfo.Operation}
	c.stateKeeper.UpdateChannel() <- update_object
}

func toReadableState(lwwes crdt.Lwwes) ReadableState {
	values := make([]string, 0)
	for _, element := range lwwes.Get() {
		intElement, ok := element.(domain.IntElement)
		if ok {
			values = append(values, strconv.Itoa(intElement.Value))
		}
	}
	result := ReadableState{Values:values}
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