package discovery

import (
	"encoding/json"
	"net/http"
	"bytes"
	"time"
	"fmt"
)

type Heartbeater interface {
	PublishHeartbeat()
	Start(d time.Duration)
}

func CreateDiscoveryHeartbeater(dc *DiscoveryClient)Heartbeater {
	client := http.Client{}
	h := discoveryHeartbeater{dc,&client}
	return &h
}

type discoveryHeartbeater struct {
	discovery_client *DiscoveryClient
	client *http.Client
}

func (heartbeater *discoveryHeartbeater) Start(d time.Duration) {
	go func() {
		for  {
			time.Sleep(d)
			heartbeater.PublishHeartbeat()
		}
	}()
}

func (heartbeater *discoveryHeartbeater) PublishHeartbeat() {
	client := heartbeater.discovery_client
	info := (*client).HeartbeatInfo()
	fmt.Println(fmt.Sprintf("publishing heartbeat to %d nodes", len(info.Cluster.Nodes)))
	for _, node := range info.Cluster.Nodes {
		heartbeater.sendHeartbeatInfo(info, node.Url)
	}
}

func (dc *discoveryHeartbeater) sendHeartbeatInfo(info HeartbeatInfo, target string) {
	jsonVal, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	rq, _ := http.NewRequest("POST", target+"/cluster/heartbeat", bytes.NewBuffer(jsonVal))
	dc.client.Do(rq)
}


