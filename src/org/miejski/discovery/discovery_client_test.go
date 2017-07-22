package discovery

import (
	"testing"
)

func TestActiveNodes(t *testing.T) {
	client := NewDiscoveryClient()
	client.AddNode(AppNode{State: ACTIVE})
	client.AddNode(AppNode{State: DEAD})
	client.AddNode(AppNode{State: ACTIVE})

	activeNodes := client.CurrentActiveNodes()

	if len(activeNodes) != 2 {
		t.Errorf("ActiveNodes count = %d", len(activeNodes))
	}
}

func TestClusterStatus(t *testing.T) {
	client := NewDiscoveryClient()
	client.RegisterHeartbeat(NodeInfo{"node1"})
	client.RegisterHeartbeat(NodeInfo{"node2"})
	client.RegisterHeartbeat(NodeInfo{"node3"})
	client.RegisterHeartbeat(NodeInfo{"node3"})
	client.RegisterHeartbeat(NodeInfo{"node1"})

	clusterInfo := client.ClusterInfo()

	if len(clusterInfo.nodes) != 3 {
		t.Errorf("Cluster info should have %d nodes, but has %d!", 3, len(clusterInfo.nodes))
	}
}

