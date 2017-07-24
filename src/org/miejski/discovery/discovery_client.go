package discovery

import (
	"time"
	"fmt"
	"net/http"
	"encoding/json"
)

type DiscoveryClient interface {
	ClusterInfo() ClusterStatus
	CurrentActiveNodes() []AppNode
	AllNodes() []AppNode
	AddNode(node AppNode)
	RegisterHeartbeat(node_info HeartbeatInfo)
	HeartbeatInfo() HeartbeatInfo
}

func NewDiscoveryClient(thisUrl string, joinAddress string) DiscoveryClient {
	nodes := make([]AppNode, 0)
	client := inMemoryDiscoveryClient{info: AppNode{Url: thisUrl}, Nodes: nodes}
	if joinAddress != "" {
		fmt.Println(fmt.Sprintf("Joining cluster at: %s", joinAddress))
		client.UpdateClusterInfo(joinAddress)
		fmt.Println(fmt.Sprintf("Joined cluster: %#v", client.ClusterInfo()))
	}
	return &client
}

type inMemoryDiscoveryClient struct {
	info  AppNode
	Nodes []AppNode
}

func (client *inMemoryDiscoveryClient) UpdateClusterInfo(joinAddress string) {
	cluster_info := getClusterInfo(joinAddress)
	nodes := make([]AppNode, 0)
	for _, node := range cluster_info.Nodes {
		nodes = append(nodes, AppNode{node.Url, ACTIVE, time.Now()})
	}
	nodes = append(nodes, AppNode{cluster_info.NodeUrl, ACTIVE, time.Now()})
	client.Nodes = nodes
}
func getClusterInfo(joinAddress string) ClusterStatus {
	rq, _ := http.NewRequest(http.MethodGet, joinAddress+"/cluster/info", nil)
	client := http.Client{}
	rs, _ := client.Do(rq)

	decoder := json.NewDecoder(rs.Body)
	var value ClusterStatus
	err := decoder.Decode(&value)
	if err != nil {
		panic(err)
	}
	return value
}

func (client *inMemoryDiscoveryClient) HeartbeatInfo() HeartbeatInfo {
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
	// TODO
}

func (client *inMemoryDiscoveryClient) AllNodes() []AppNode {
	return client.Nodes
}
