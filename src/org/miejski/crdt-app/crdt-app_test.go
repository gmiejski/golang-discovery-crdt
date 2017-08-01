package main

import (
	"testing"
	"os"
	"flag"
	"org/miejski/discovery"
	"org/miejski/domain"
	"net/http"
	"fmt"
	"encoding/json"
	"org/miejski/crdt"
	"time"
	"bytes"
	"github.com/stretchr/testify/assert"
	"org/miejski/simple_json"
)
var port int = 7778
var host string = fmt.Sprintf("http://localhost:%d", port)

func setup() {
	reset(host)
}

func TestMain(m *testing.M) {
	result := 0
	flag.Parse()

	if !testing.Short() {
		server := prepareServer()
		result = m.Run()
		quitServer(server)
	}
	os.Exit(result)
}

func prepareServer() CrdtServer {
	discovery_client := discovery.NewDiscoveryClient("host", "")
	unsafe_keeper := domain.UnsafeDomainKeeper()
	dk := CreateSafeValueKeeper(&unsafe_keeper)
	state_controller := newStateController(&discovery_client, &dk)
	server := NewServer(&state_controller, &discovery_client)

	go func() {
		server.Start(port)
	}()

	return server
}

func quitServer(server CrdtServer) {
	server.Shutdown()
}

func TestCrdtUpdateIntegration(t *testing.T) {
	setup()
	// given
	first_expected := crdt.CreateLwwes()
	now := time.Now()
	first_expected.Add(&domain.IntElement{1}, now)
	second_expected := crdt.CreateLwwes()
	second_expected.Add(&domain.IntElement{1}, now)
	second_expected.Add(&domain.IntElement{2}, now)

	updateValue(1, domain.ADD, host)
	updateValue(2, domain.ADD, host)

	// when
	first_value := getCurrentValue(host)

	// then
	assert.True(t, first_value.Contains(domain.IntElement{1}))
	assert.True(t, first_value.Contains(domain.IntElement{2}))

	// when
	updateValue(2, domain.REMOVE, host)
	second_value := getCurrentValue(host)

	// then
	assert.True(t, second_value.Contains(domain.IntElement{1}))
	assert.False(t, second_value.Contains(domain.IntElement{2}))
}

func TestCrdtResetIntegration(t *testing.T) {
	setup()
	// given
	first_expected := crdt.CreateLwwes()
	now := time.Now()
	first_expected.Add(&domain.IntElement{1}, now)
	second_expected := crdt.CreateLwwes()
	second_expected.Add(&domain.IntElement{1}, now)
	second_expected.Add(&domain.IntElement{2}, now)
	updateValue(1, domain.ADD, host)
	updateValue(2, domain.ADD, host)

	// when
	first_value := getCurrentValue(host)

	// then
	assert.True(t, first_value.Contains(domain.IntElement{1}))
	assert.True(t, first_value.Contains(domain.IntElement{2}))

	// when
	reset(host)
	second_value := getCurrentValue(host)

	// then
	assert.False(t, second_value.Contains(domain.IntElement{1}))
	assert.False(t, second_value.Contains(domain.IntElement{2}))
}

func getCurrentValue(host string) crdt.Lwwes {
	rq, _ := http.NewRequest("GET", host+"/status", nil)
	client := http.Client{}
	rs, _ := client.Do(rq)
	var dv CurrentStateDto
	simple_json.Unmarshal(rs.Body, &dv)
	return lwwesFromDto(dv)
}

func updateValue(val int, op domain.UpdateOperationType, host string) {
	client := http.Client{}
	operation := CrdtOperation{domain.IntElement{Value: val}, op}
	jsonVal, err := json.Marshal(operation)
	if err != nil {
		panic(err)
	}
	body := bytes.NewBuffer(jsonVal)
	rq, _ := http.NewRequest("POST", host+"/status/update", body)
	client.Do(rq)
}

func reset(host string) {
	rq, _ := http.NewRequest("POST", host+"/status/reset", nil)
	client := http.Client{}
	client.Do(rq)
}

func TestElementsMap2(t *testing.T) {
	simple_map := map[crdt.Element]time.Time{}
	impl := domain.IntElement{1}
	simple_map[impl] = time.Now()

	impl2 := domain.IntElement{2}
	simple_map[impl2] = time.Now()

	aad := domain.IntElement{1}
	_, v2 := simple_map[aad]
	//print(v1)
	if !v2 {
		t.Fail()
	}
}
