package discovery

type DiscoveryClient interface {
	CurrentActiveNodes() []AppNode
	AllNodes() []AppNode
	AddNode(node AppNode)
}

type inMemoryDiscoveryClient struct {
	Nodes []AppNode
}

func NewDiscoveryClient() DiscoveryClient {
	nodes := make([]AppNode, 5)
	client := inMemoryDiscoveryClient{Nodes: nodes}
	return &client
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
