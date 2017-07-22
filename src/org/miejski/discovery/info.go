package discovery

type NodeInfo struct {
	url string
}

type HeartbeatInfo struct {
	url string
	cluster ClusterStatus
}