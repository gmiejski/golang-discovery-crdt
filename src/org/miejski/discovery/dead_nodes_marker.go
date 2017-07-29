package discovery

import (
	"time"
	"fmt"
)

var time_alive time.Duration = 5 * time.Second

type DeadNodeMarker interface {
	StartMarking(every time.Duration)
}

type deadNodeMarker struct {
	dc *DiscoveryClient
}

func (marker *deadNodeMarker) StartMarking(every time.Duration) {
	go start(marker.dc, every)
}

func NewDeadNodeMarker(dc *DiscoveryClient) DeadNodeMarker {
	marker := deadNodeMarker{dc}
	return &marker
}

func start(dc *DiscoveryClient, every time.Duration) {
	for {
		time.Sleep(every)
		markDeadNodes(dc)
	}
}

func markDeadNodes(dc *DiscoveryClient) {
	fmt.Println("Start marking nodes as DEAD")
	dead_nodes := make([]AppNode, 0)
	for _, node := range (*dc).AllNodes() {
		if node.State != DEAD && node.LastUpdate.Add(time_alive).Before(time.Now()) {
			dead_nodes = append(dead_nodes, node)
			(*dc).ChangeStatus(node.Url, node.LastUpdate, DEAD)
		}
	}
	fmt.Println(fmt.Sprintf("Nodes marked as dead #%v", dead_nodes))
}
