package discovery

import (
	"testing"
	"net/http/httptest"
	"net/http"
)

var heartbeatsReceived = 0

func TestHeartbeatingOtherNodesInCluster(t *testing.T) {
	// given
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		heartbeatsReceived = heartbeatsReceived + 1
	}))
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		heartbeatsReceived = heartbeatsReceived + 1
	}))
	defer ts.Close()
	defer ts2.Close()
	discovery_client := NewDiscoveryClient("http://localhost:8080")
	discovery_client.RegisterHeartbeat(HeartbeatInfo{Url: ts.URL})
	discovery_client.RegisterHeartbeat(HeartbeatInfo{Url: ts2.URL})
	heartbeater := CreateDiscoveryHeartbeater(&discovery_client)

	// when
	heartbeater.PublishHeartbeat()

	// then

	if heartbeatsReceived != 2 {
		t.Errorf("Bad number of requests send: expected %d got %d!", 2, heartbeatsReceived)
	}
}
