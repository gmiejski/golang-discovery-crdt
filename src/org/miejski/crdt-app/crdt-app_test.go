package main

import (
	"testing"
	"os"
	"flag"
	"org/miejski/discovery"
	"org/miejski/domain"
	"net/http"
	"io/ioutil"
	"strconv"
)

var host string = "http://localhost:8080"

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
	discovery_client := discovery.NewDiscoveryClient()
	dk := domain.UnsafeDomainKeeper()
	state_controller := newStateController(&discovery_client, &dk)
	server := NewServer(state_controller)

	go func() {
		server.Start(8080)
	}()

	return server
}

func quitServer(server CrdtServer) {
	server.Shutdown()
}

func TestCrdtUpdateIntegration(t *testing.T) {
	setup()
	// given
	first_expected := domain.DomainValue(0)
	second_expected := domain.DomainValue(2)

	// when
	first_value := getCurrentValue(host)

	// then
	if first_value != domain.DomainValue(first_expected) {
		t.Errorf("Should have value %d, but had %d before updates", first_expected, first_value)
	}

	// when
	updateValue(host)
	updateValue(host)
	second_value := getCurrentValue(host)

	// then
	if second_value != domain.DomainValue(second_expected) {
		t.Errorf("Should have value %d, but had %d after updates", second_expected, second_value)
	}
}

func TestCrdtResetIntegration(t *testing.T) {
	setup()
	// given
	first_expected := domain.DomainValue(2)
	after_reset_expected := domain.DomainValue(0)

	// when
	updateValue(host)
	updateValue(host)
	first_value := getCurrentValue(host)

	// then
	if first_value != domain.DomainValue(first_expected) {
		t.Errorf("Should have value %d, but had %d before updates", first_expected, first_value)
	}

	// when
	reset(host)
	second_value := getCurrentValue(host)

	// then
	if second_value != domain.DomainValue(after_reset_expected) {
		t.Errorf("Should have value %d, but had %d after updates", after_reset_expected, second_value)
	}
}

func getCurrentValue(host string) domain.DomainValue {
	rq, _ := http.NewRequest("GET", host+"/status", nil)
	client := http.Client{}
	rs, _ := client.Do(rq)
	bodyBytes, _ := ioutil.ReadAll(rs.Body)
	domain_value, _ := strconv.Atoi(string(bodyBytes))
	return domain.DomainValue(domain_value)
}

func updateValue(host string) {
	rq, _ := http.NewRequest("POST", host+"/status/update", nil)
	client := http.Client{}
	client.Do(rq)
}

func reset(host string) {
	rq, _ := http.NewRequest("POST", host+"/status/reset", nil)
	client := http.Client{}
	client.Do(rq)
}
