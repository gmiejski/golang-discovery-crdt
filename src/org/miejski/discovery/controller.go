package discovery

import (
	"net/http"
	"org/miejski/rest"
	"fmt"
)

type DiscoveryController interface {
	ClusterInfo(w http.ResponseWriter, request *http.Request)
	Heartbeat(w http.ResponseWriter, request *http.Request)
}

type simpleDiscoveryController struct{
	discovery_client *DiscoveryClient
}

func (dc *simpleDiscoveryController) ClusterInfo(w http.ResponseWriter, request *http.Request) {
	fmt.Fprint(w, "ClusterInfo")
}

func (dc *simpleDiscoveryController) Heartbeat(w http.ResponseWriter, request *http.Request) {
	fmt.Fprint(w, "Heartbeat")
}

func CreateDiscoveryController(discovery_client *DiscoveryClient) DiscoveryController {
	dc := simpleDiscoveryController{discovery_client}
	return &dc
}

func RegisterDiscoveryEndpoints(discovery_client *DiscoveryClient) {
	dc := CreateDiscoveryController(discovery_client)
	http.HandleFunc("/cluster/info", rest.GET(dc.ClusterInfo))
	http.HandleFunc("/cluster/heartbeat", rest.POST(dc.Heartbeat))
}

