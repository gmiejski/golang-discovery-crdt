package discovery

import (
	"testing"
	"net/http"
	"github.com/stretchr/testify/mock"
	"fmt"
	"encoding/json"
	"bytes"
	"time"
)

var client = http.Client{}
var port int = 7779
var url string = fmt.Sprintf("http://localhost:%d", port)

func TestControllerUnmarshallingHeartbeat(t *testing.T) {
	// given
	dc := new(MockDiscoveryClient)
	c := DiscoveryClient(dc)
	mux := http.NewServeMux()
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler:mux}
	expected_info := HeartbeatInfo{Url: url}
	RegisterDiscoveryEndpoints(&c, mux)
	go startServer(srv, &c, t)

	dc.On("RegisterHeartbeat", expected_info)

	// when
	client.Do(createRequestWithHearbeat(url, expected_info))

	// then
	dc.AssertCalled(t, "RegisterHeartbeat", expected_info)
	defer srv.Shutdown(nil)
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
func startServer(srv *http.Server, dc *DiscoveryClient, t *testing.T) {
	//RegisterDiscoveryEndpoints(dc, srv)
	//CreateDiscoveryHeartbeater(dc).Start(2 * time.Second)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

// MOCKS

type MockDiscoveryClient struct {
	mock.Mock
}

func (c *MockDiscoveryClient) ChangeStatus(node_url string, last_time_seen time.Time, new_state State) {
	panic("implement me")
}

func (c *MockDiscoveryClient) Start(every time.Duration) {
	panic("implement me")
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
