package discovery

import (
	"time"
	"io"
	"encoding/json"
)

type State string

type AppNode struct {
	Url        string
	State      State
	LastUpdate time.Time
}

const (
	ACTIVE State = "ACTIVE"
	DEAD   State = "DEAD"
)

type ClusterStatus struct {
	NodeUrl string
	Nodes   []AppNode
}

func (c *ClusterStatus) Unmarshal(data io.ReadCloser) error {
	decoder := json.NewDecoder(data)
	err := decoder.Decode(c)
	if err != nil {
		panic(err)
	}
	return nil
}
