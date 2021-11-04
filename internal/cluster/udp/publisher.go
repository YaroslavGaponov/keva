package udp

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/YaroslavGaponov/keva/pkg/cluster"
	"github.com/YaroslavGaponov/keva/pkg/utils"
)

type UdpPublisher struct {
	options UdpOptions
}

func NewPublisher(options UdpOptions) *UdpPublisher {
	c := UdpPublisher{
		options: options,
	}
	return &c
}

func (c *UdpPublisher) RegisterNode(info cluster.NodeInfo) error {

	ServerAddr, err := net.ResolveUDPAddr("udp", c.options.BroadcastHost+":"+c.options.BroadcastPort)
	if err != nil {
		return err
	}

	ip := utils.SelfIP()
	LocalAddr, err := net.ResolveUDPAddr("udp", ip.String()+":0")
	if err != nil {
		return err
	}

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err != nil {
		return err
	}

	ping := Message{
		Cmd:     CMD_PING,
		Cluster: c.options.ClusterName,
		Info:    info,
	}

	raw, err := json.Marshal(ping)
	if err != nil {
		return err
	}

	go func() {
		for {
			if _, err := Conn.Write(raw); err != nil {
				log.Println(err)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	return nil
}
