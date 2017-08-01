package discovery

import (
	"net/http"
	"org/miejski/rest"
	"fmt"
	"encoding/json"
	"org/miejski/simple_json"
)

type DiscoveryController interface {
	ClusterInfo(w http.ResponseWriter, request *http.Request)
	Heartbeat(w http.ResponseWriter, request *http.Request)
}

type simpleDiscoveryController struct{
	discovery_client *DiscoveryClient
	client *http.Client
}

func (dc *simpleDiscoveryController) ClusterInfo(w http.ResponseWriter, request *http.Request) {
	client := dc.discovery_client
	info := (*client).ClusterInfo()
	val, err := json.Marshal(&info)
	if err != nil {
		panic(err)
	}
 	fmt.Fprint(w, string(val))
}

func (dc *simpleDiscoveryController) Heartbeat(w http.ResponseWriter, request *http.Request) {
	heartbeat_info := readHeartbeatInfo(request)
	ds := dc.discovery_client
	(*ds).RegisterHeartbeat(heartbeat_info)
}

func readHeartbeatInfo(request *http.Request) HeartbeatInfo {
	var t HeartbeatInfo
	simple_json.Unmarshal(request.Body, &t)
	defer request.Body.Close()
	return t
}

func CreateDiscoveryController(discovery_client *DiscoveryClient) DiscoveryController {
	client := http.Client{}
	dc := simpleDiscoveryController{discovery_client, &client}
	return &dc
}

func RegisterDiscoveryEndpoints(discovery_client *DiscoveryClient, srv *http.ServeMux) {
	dc := CreateDiscoveryController(discovery_client)
	srv.HandleFunc("/cluster/info", rest.GET(dc.ClusterInfo))
	srv.HandleFunc("/cluster/heartbeat", rest.POST(dc.Heartbeat))
}

