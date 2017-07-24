package discovery

import (
	"time"
	"fmt"
)

type DiscoveryClient interface {
	ClusterInfo() ClusterStatus
	CurrentActiveNodes() []AppNode
	AllNodes() []AppNode
	AddNode(node AppNode)
	RegisterHeartbeat(node_info HeartbeatInfo)
	HeartbeatInfo() HeartbeatInfo
}

func NewDiscoveryClient(thisUrl string) DiscoveryClient {
	nodes := make([]AppNode, 0)
	client := inMemoryDiscoveryClient{info: AppNode{Url: thisUrl}, Nodes: nodes}
	return &client
}

type inMemoryDiscoveryClient struct {
	info  AppNode
	Nodes []AppNode
}

func (client *inMemoryDiscoveryClient) HeartbeatInfo() HeartbeatInfo {
	fmt.Println(client.info)
	return HeartbeatInfo{client.info.Url, client.ClusterInfo()}
}

func (client *inMemoryDiscoveryClient) RegisterHeartbeat(node_info HeartbeatInfo) {
	exists, node := client.containsNode(node_info.Url)
	if exists {
		node.State = ACTIVE
		node.LastUpdate = time.Now()
	} else {
		node := AppNode{Url: node_info.Url, State: ACTIVE, LastUpdate: time.Now()}
		client.AddNode(node)
	}
}
func (client *inMemoryDiscoveryClient) containsNode(url string) (bool, *AppNode) {
	for _, v := range client.Nodes {
		if v.Url == url {
			return true, &v
		}
	}
	return false, nil
}

func (client *inMemoryDiscoveryClient) ClusterInfo() ClusterStatus {
	return ClusterStatus{client.info.Url, client.Nodes}
}

func (client *inMemoryDiscoveryClient) CurrentActiveNodes() []AppNode {
	result := make([]AppNode, 0)
	for _, v := range client.AllNodes() {
		if v.State == ACTIVE {
			result = append(result, v)
		}
	}
	return result
}

func (client *inMemoryDiscoveryClient) AddNode(node AppNode) {
	client.Nodes = append(client.Nodes, node)
}

func (client *inMemoryDiscoveryClient) UpdateNodeState(node AppNode) {

}

func (client *inMemoryDiscoveryClient) AllNodes() []AppNode {
	return client.Nodes
}
