package discovery

import "time"

type DiscoveryClient interface {
	ClusterInfo() ClusterStatus
	CurrentActiveNodes() []AppNode
	AllNodes() []AppNode
	AddNode(node AppNode)
	RegisterHeartbeat(node_info NodeInfo)
}

type inMemoryDiscoveryClient struct {
	Nodes []AppNode
}

func NewDiscoveryClient() DiscoveryClient {
	nodes := make([]AppNode, 0)
	client := inMemoryDiscoveryClient{Nodes: nodes}
	return &client
}

func (client *inMemoryDiscoveryClient) RegisterHeartbeat(node_info NodeInfo) {
	exists, node := client.containsNode(node_info.url)
	if exists {
		node.State = ACTIVE
		node.LastUpdate = time.Now()
	} else {
		node := AppNode{url:node_info.url, State:ACTIVE, LastUpdate:time.Now()}
		client.AddNode(node)
	}
}
func (client *inMemoryDiscoveryClient) containsNode(url string) (bool, *AppNode) {
	for _, v := range client.Nodes {
		if v.url == url {
			return true, &v
		}
	}
	return false, nil
}

func (client *inMemoryDiscoveryClient) ClusterInfo() ClusterStatus {
	return ClusterStatus{client.Nodes}
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
