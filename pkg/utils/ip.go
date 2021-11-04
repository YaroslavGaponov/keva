package utils

import (
	"net"
)

func SelfIP() net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP
			}
		}
	}

	return net.ParseIP("127.0.0.1")
}

func GetTcpListener(host string) (net.Listener, error) {
	listener, err := net.Listen("tcp", host+":0")
	if err != nil {
		return nil, err
	}
	return listener, nil

}
