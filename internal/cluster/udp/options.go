package udp

import (
	"github.com/YaroslavGaponov/keva/pkg/utils"
)

type UdpOptions struct {
	ClusterName   string
	BroadcastHost string
	BroadcastPort string
}

func OptionsFromEnv() UdpOptions {
	clusterName := utils.GetEnvVariableOrDefult("CLUSTER_NAME", "keva")
	broadcasthost := utils.GetEnvVariableOrDefult("BROADCAST_HOST", "224.2.2.4")
	broadcastport := utils.GetEnvVariableOrDefult("BROADCAST_PORT", "8787")
	return UdpOptions{
		ClusterName:   clusterName,
		BroadcastHost: broadcasthost,
		BroadcastPort: broadcastport,
	}
}
