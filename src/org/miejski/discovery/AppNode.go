package discovery

import "time"

type State string

type AppNode struct {
	Url        string
	State      State
	LastUpdate time.Time
}

const (
	ACTIVE State = "ACTIVE"
	DEAD State = "DEAD"
)

type ClusterStatus struct {
	NodeUrl string
	Nodes   []AppNode
}