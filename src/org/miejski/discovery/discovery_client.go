package discovery

import (
	"time"
	"fmt"
	"net/http"
	"encoding/json"
	"sync"
)

var time_alive time.Duration = 5 * time.Second

type DiscoveryClient interface {
	ClusterInfo() ClusterStatus
	CurrentActiveNodes() []AppNode
	AllNodes() []AppNode
	RegisterHeartbeat(node_info HeartbeatInfo)
	HeartbeatInfo() HeartbeatInfo
	Start(every time.Duration)
	ChangeStatus(node_url string, last_time_seen time.Time, new_state State)
}

func NewDiscoveryClient(thisUrl string, joinAddress string) DiscoveryClient {
	nodes := make([]AppNode, 0)
	client := inMemoryDiscoveryClient{info: AppNode{Url: thisUrl}, Nodes: nodes}
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

func (dc *inMemoryDiscoveryClient) Start(every time.Duration) {
	go start(dc, every)
}

func start(dc *inMemoryDiscoveryClient, every time.Duration) {
	for {
		time.Sleep(every)
		dc.l.Lock()
		markDeadNodes(dc)
		dc.l.Unlock()
	}
}
func markDeadNodes(dc *inMemoryDiscoveryClient) {
	fmt.Println("Start marking nodes as DEAD")
	dead_nodes := make([]AppNode, 0)
	for _, node := range (*dc).AllNodes() {
		if node.State != DEAD && node.LastUpdate.Add(time_alive).Before(time.Now()) {
			dead_nodes = append(dead_nodes, node)
			dc.ChangeStatus(node.Url, node.LastUpdate, DEAD)
		}
	}
	fmt.Println(fmt.Sprintf("Nodes marked as dead #%v", dead_nodes))
}

func (client *inMemoryDiscoveryClient) updateClusterInfo(joinAddress string) {
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
	client.l.Lock()
	defer client.l.Unlock()
	exists, node := client.containsNode(node_info.Url)
	if exists {
		fmt.Println(fmt.Sprintf("Marking node as ALIVE: %s", node.Url))
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
	client.l.Lock()
	defer client.l.Unlock()
	result := make([]AppNode, 0)
	for _, v := range client.AllNodes() {
		if v.State == ACTIVE {
			result = append(result, v)
		}
	}
	return result
}

func (client *inMemoryDiscoveryClient) addNode(node AppNode) { // TODO remove
	client.Nodes = append(client.Nodes, node)
}

func (client *inMemoryDiscoveryClient) AllNodes() []AppNode { // TODO remove?
	return client.Nodes
}
