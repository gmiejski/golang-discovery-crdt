package discovery

import (
	"time"
	"fmt"
	"net/http"
	"sync"
	"org/miejski/simple_json"
)

type DiscoveryClient interface {
	ClusterInfo() ClusterStatus
	CurrentActiveNodes() []AppNode
	AllNodes() []AppNode
	RegisterHeartbeat(node_info HeartbeatInfo)
	HeartbeatInfo() HeartbeatInfo
	ChangeStatus(node_url string, last_time_seen time.Time, new_state State)
}

func NewDiscoveryClient(thisUrl string, joinAddress string) DiscoveryClient {
	nodes := make([]AppNode, 0)
	http_client := http.Client{
		Timeout: time.Second * 4,
	}
	client := inMemoryDiscoveryClient{info: AppNode{Url: thisUrl}, Nodes: nodes, client: &http_client}
	if joinAddress != "" {
		fmt.Println(fmt.Sprintf("Joining cluster at: %s", joinAddress))
		client.updateClusterInfo(joinAddress)
		fmt.Println(fmt.Sprintf("Joined cluster: %#v", client.ClusterInfo()))
	}
	return &client
}

type inMemoryDiscoveryClient struct {
	info  AppNode
	Nodes []AppNode
	l     sync.Mutex
	client *http.Client
}

func (dc *inMemoryDiscoveryClient) ChangeStatus(node_url string, last_time_seen time.Time, new_state State) {
	for i, a := range dc.Nodes {
		if a.Url == node_url && a.LastUpdate == last_time_seen {
			dc.Nodes[i].State = new_state
			dc.Nodes[i].LastUpdate = time.Now()
			break
		}
	}
}

func (client *inMemoryDiscoveryClient) updateClusterInfo(joinAddress string) {
	cluster_info := client.getClusterInfo(joinAddress)
	nodes := make([]AppNode, 0)
	for _, node := range cluster_info.Nodes {
		if node.Url != client.info.Url {
			nodes = append(nodes, AppNode{node.Url, node.State, time.Now()})
		}
	}
	nodes = append(nodes, AppNode{cluster_info.NodeUrl, ACTIVE, time.Now()})
	client.Nodes = nodes
}

func (client *inMemoryDiscoveryClient) getClusterInfo(joinAddress string) ClusterStatus {
	rq, _ := http.NewRequest(http.MethodGet, joinAddress+"/cluster/info", nil)

	rs, _ := client.client.Do(rq)

	var value ClusterStatus
	simple_json.Unmarshal(rs.Body, &value)
	return value
}

func (client *inMemoryDiscoveryClient) HeartbeatInfo() HeartbeatInfo {
	return HeartbeatInfo{client.info.Url, client.ClusterInfo()}
}

func (client *inMemoryDiscoveryClient) RegisterHeartbeat(node_info HeartbeatInfo) {
	client.l.Lock()
	defer client.l.Unlock()
	exists, node := client.containsNode(node_info.Url)
	if exists {
		//fmt.Println(fmt.Sprintf("Marking node as ALIVE: %s", node.Url))
		client.ChangeStatus(node.Url, node.LastUpdate, ACTIVE)
	} else {
		node := AppNode{Url: node_info.Url, State: ACTIVE, LastUpdate: time.Now()}
		client.addNode(node)
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
	client.l.Lock()
	defer client.l.Unlock()
	status := ClusterStatus{client.info.Url, client.Nodes}
	return status
}

func (client *inMemoryDiscoveryClient) CurrentActiveNodes() []AppNode {
	nodes := client.AllNodes()
	client.l.Lock()
	defer client.l.Unlock()
	result := make([]AppNode, 0)
	for _, v := range nodes {
		if v.State == ACTIVE {
			result = append(result, v)
		}
	}
	return result
}

func (client *inMemoryDiscoveryClient) addNode(node AppNode) { // TODO remove
	client.Nodes = append(client.Nodes, node)
}

func (client *inMemoryDiscoveryClient) AllNodes() []AppNode {
	client.l.Lock()
	defer client.l.Unlock()
	nodes := client.Nodes
	return nodes
}
