package discovery

import (
	"testing"
	"time"
)

func TestRemarkNodeAsAliveAfterDead(t *testing.T) {
	client := NewDiscoveryClient("host", "")
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node1"})
	status := client.ClusterInfo().Nodes
	for _,n := range status {
		client.ChangeStatus(n.Url, n.LastUpdate, DEAD)
	}
	// then
	dead_nodes := dead_nodes_count(&client)
	if dead_nodes != 1 {
		t.Errorf("Cluster info should have %d DEAD Nodes, but has %d!", 1, dead_nodes)
	}

	// when
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node1"})

	// then
	dead_nodes = dead_nodes_count(&client)
	if dead_nodes != 0 {
		t.Errorf("Cluster info should have %d DEAD Nodes, but has %d!", 0, dead_nodes)
	}
}

func TestMarkingNodeAsDead(t *testing.T) {
	// given
	time_alive = 1 * time.Millisecond

	client := NewDiscoveryClient("host", "")
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node1"})
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node2"})

	time.Sleep(100 * time.Millisecond)

	// when
	markDeadNodes(&client)

	// then
	dead_nodes := dead_nodes_count(&client)
	if dead_nodes != 2 {
		t.Errorf("Cluster info should have %d DEAD Nodes, but has %d!", 2, dead_nodes)
	}
}
