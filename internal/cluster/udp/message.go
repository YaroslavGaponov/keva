package udp

import "github.com/YaroslavGaponov/keva/pkg/cluster"

const (
	CMD_PING = "ping"
)

type Message struct {
	Cmd     string           `json:"cmd"`
	Cluster string           `json:"cluster"`
	Info    cluster.NodeInfo `json:"info"`
}
