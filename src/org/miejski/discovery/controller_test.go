package discovery

import (
	"testing"
	"net/http"
	"github.com/stretchr/testify/mock"
	"fmt"
	"encoding/json"
	"bytes"
)

var client = http.Client{}
var url string = "http://localhost:8080"

func TestControllerUnmarshallingHeartbeat(t *testing.T) {
	// given
	dc := new(MockDiscoveryClient)
	c := DiscoveryClient(dc)
	srv := &http.Server{Addr: fmt.Sprintf(":%d", 8080)}
	RegisterDiscoveryEndpoints(&c, srv)
	expected_info := HeartbeatInfo{Url: url}
	go startServer(srv)
	defer srv.Shutdown(nil)
	dc.On("RegisterHeartbeat", expected_info)

	// when
	client.Do(createRequestWithHearbeat(url, expected_info))

	// then
	dc.AssertCalled(t, "RegisterHeartbeat", expected_info)

}
func createRequestWithHearbeat(url string, info HeartbeatInfo) (*http.Request) {
	jsonVal, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	body := bytes.NewBuffer(jsonVal)
	rq, _ := http.NewRequest("POST", url+"/cluster/heartbeat", body)
	return rq
}
func startServer(srv *http.Server) {
	srv.ListenAndServe()
}

// MOCKS

type MockDiscoveryClient struct {
	mock.Mock
}

func (c *MockDiscoveryClient) ClusterInfo() ClusterStatus {
	panic("implement me")
}

func (c *MockDiscoveryClient) CurrentActiveNodes() []AppNode {
	panic("implement me")
}

func (c *MockDiscoveryClient) AllNodes() []AppNode {
	panic("implement me")
}

func (c *MockDiscoveryClient) AddNode(node AppNode) {
	panic("implement me")
}

func (m *MockDiscoveryClient) RegisterHeartbeat(node_info HeartbeatInfo) {
	m.Called(node_info)
}

func (c *MockDiscoveryClient) HeartbeatInfo() HeartbeatInfo {
	panic("implement me")
}
