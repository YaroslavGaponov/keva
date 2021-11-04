package cluster

type NodeInfo struct {
	NodeId string `json:"nodeId"`
	Schema string `json:"schema"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Path   string `json:"path"`
}
