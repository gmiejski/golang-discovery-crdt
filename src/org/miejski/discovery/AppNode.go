package discovery

import "time"

type State int

type AppNode struct {
	url string
	State State
	LastUpdate time.Time
}

const (
	DEAD State = iota
	ACTIVE State = iota
)

type ClusterStatus struct {
	nodes []AppNode
}