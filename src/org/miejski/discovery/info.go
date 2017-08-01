package discovery

import (
	"io"
	"encoding/json"
)

type NodeInfo struct {
	Url string
}

type HeartbeatInfo struct {
	Url     string
	Cluster ClusterStatus
}

func (h *HeartbeatInfo) Unmarshal(data io.ReadCloser) error {
	decoder := json.NewDecoder(data)
	defer data.Close()
	err := decoder.Decode(h)
	if err != nil {
		return err
	}
	return nil

}
