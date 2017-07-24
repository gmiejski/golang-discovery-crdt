package discovery

type NodeInfo struct {
	Url string
}

type HeartbeatInfo struct {
	Url     string
	Cluster ClusterStatus
}