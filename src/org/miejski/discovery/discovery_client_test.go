package discovery

import (
	"testing"
)

func TestActiveNodes(t *testing.T) {
	client := NewDiscoveryClient()
	client.AddNode(AppNode{State:ACTIVE})
	client.AddNode(AppNode{State:DEAD})
	client.AddNode(AppNode{State:ACTIVE})

	activeNodes := client.CurrentActiveNodes()

	if len(activeNodes) != 2 {
		t.Errorf("ActiveNodes count = %d", len(activeNodes))
	}
}

