package discovery

import (
	"testing"
)

//func TestActiveNodes(t *testing.T) {
//	client := NewDiscoveryClient("host", "")
//	client.AddNode(AppNode{State: ACTIVE})
//	client.AddNode(AppNode{State: DEAD})
//	client.AddNode(AppNode{State: ACTIVE})
//
//	activeNodes := client.CurrentActiveNodes()
//
//	if len(activeNodes) != 2 {
//		t.Errorf("ActiveNodes count = %d", len(activeNodes))
//	}
//}

func TestClusterStatus(t *testing.T) {
	client := NewDiscoveryClient("host", "")
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node1"})
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node2"})
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node3"})
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node3"})
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node1"})

	clusterInfo := client.ClusterInfo()

	if len(clusterInfo.Nodes) != 3 {
		t.Errorf("Cluster info should have %d Nodes, but has %d!", 3, len(clusterInfo.Nodes))
	}
}

func dead_nodes_count(client *DiscoveryClient) int {
	result := 0
	for _, x:= range (*client).ClusterInfo().Nodes {
		if x.State == DEAD {
			result ++
		}
	}
	return result
}

func TestChangeStatus(t *testing.T) {
	client := NewDiscoveryClient("host", "")
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node1"})
	client.RegisterHeartbeat(HeartbeatInfo{Url: "node2"})

	status := client.ClusterInfo().Nodes
	for _,n := range status {
		client.ChangeStatus(n.Url, n.LastUpdate, DEAD)
	}

	dead_nodes := dead_nodes_count(&client)
	if dead_nodes != 2 {
		t.Errorf("Cluster info should have %d DEAD Nodes, but has %d!", 2, dead_nodes)
	}
}