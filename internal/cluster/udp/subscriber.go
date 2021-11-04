package udp

import (
	"encoding/json"
	"net"
	"time"

	"github.com/YaroslavGaponov/keva/pkg/cluster"
	"github.com/YaroslavGaponov/keva/pkg/logger"
)

type UdpSubscliber struct {
	options UdpOptions

	logger *logger.Logger

	buffer chan Message

	nodes        map[string]cluster.NodeInfo
	nodes_health map[string]time.Time
}

func NewSubscreber(options UdpOptions) *UdpSubscliber {
	c := UdpSubscliber{
		options: options,

		logger: logger.CreateLogger(),

		buffer: make(chan Message),

		nodes:        make(map[string]cluster.NodeInfo),
		nodes_health: make(map[string]time.Time),
	}
	return &c
}

func (c *UdpSubscliber) Start() error {

	ServerAddr, err := net.ResolveUDPAddr("udp", c.options.BroadcastHost+":"+c.options.BroadcastPort)
	if err != nil {
		return err
	}
	ServerConn, err := net.ListenMulticastUDP("udp", nil, ServerAddr)
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)

	go c.onNewMessage()

	go func() {
		for {
			if n, _, err := ServerConn.ReadFromUDP(buf); err == nil {
				var message Message
				if err := json.Unmarshal(buf[0:n], &message); err == nil {
					c.buffer <- message
				}
			}
		}
	}()

	return nil
}

func (c *UdpSubscliber) GetNodes() (map[string]cluster.NodeInfo, error) {
	return c.nodes, nil
}

func (c *UdpSubscliber) onNewMessage() {
	for {
		msg := <-c.buffer

		switch msg.Cmd {
		case CMD_PING:
			if msg.Cluster == c.options.ClusterName {
				if _, ok := c.nodes[msg.Info.NodeId]; !ok {
					c.nodes[msg.Info.NodeId] = msg.Info
					c.nodes_health[msg.Info.NodeId] = time.Now()
					c.logger.Info("node %s is registered", msg.Info.NodeId)
				} else {
					c.nodes_health[msg.Info.NodeId] = time.Now()
				}
			}
		}

		for nodeId, last := range c.nodes_health {
			if last.Add(10 * time.Second).Before(time.Now()) {
				delete(c.nodes, nodeId)
				delete(c.nodes_health, nodeId)
				c.logger.Info("node %s is unregistered", nodeId)
			}
		}
	}
}
