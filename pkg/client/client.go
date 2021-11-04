package client

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/YaroslavGaponov/keva/pkg/cluster"
	"github.com/YaroslavGaponov/keva/pkg/logger"
)

var (
	ErrNoQuorum = errors.New("quorum is not reached")
	ErrTimeout  = errors.New("timeout")
)

type Client struct {
	mu         sync.Mutex
	options    *ClientOptions
	subscriber cluster.Subscriber
	logger     *logger.Logger
	compressor Compressor
}

type hash_node struct {
	hash []byte
	node string
}

func NewClient(options *ClientOptions, subscriber cluster.Subscriber) *Client {
	return &Client{
		options:    options,
		subscriber: subscriber,
		logger:     logger.CreateLogger(),
		compressor: NewLzwCompressor(),
	}
}

func (c *Client) Set(key string, value []byte) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	beforeSize := len(value)
	value, err := c.compressor.Compress(value)
	if err != nil {
		c.logger.Error("client: key %s is set with error %s", key, err)
		return false, err
	}
	afterSize := len(value)
	c.logger.Trace("client: compress %d bytes to %d bytes", beforeSize, afterSize)

	nodes, err := c.subscriber.GetNodes()
	if err != nil {
		c.logger.Error("client: set key %s with error %s", key, err)
		return false, err
	}

	result := make(chan bool)
	total := len(nodes)
	quorum := (total >> 1) + 1
	ok, bad := 0, 0

	for _, info := range nodes {
		go func(info cluster.NodeInfo) {
			client := NewNodeClient(info)
			if status, err := client.Set(key, value); err == nil && status {
				result <- true
			} else {
				result <- false
			}
		}(info)
	}

	for {
		select {
		case status := <-result:
			if status {
				ok++
			} else {
				bad++
			}
			if ok >= quorum {
				c.logger.Trace("client: key %s is set with success", key)
				return true, nil
			}
			if bad >= quorum {
				c.logger.Trace("client: key %s is set with error", key)
				return false, ErrNoQuorum
			}
			if ok+bad >= total {
				c.logger.Trace("client: key %s is set with error", key)
				return false, ErrNoQuorum
			}
		case <-time.After(c.options.timeout):
			c.logger.Trace("client: key %s is setted with timeout", key)
			return false, ErrTimeout
		}
	}
}

func (c *Client) Get(key string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	nodes, err := c.subscriber.GetNodes()
	if err != nil {
		return nil, err
	}

	result := make(chan hash_node)
	total := len(nodes)
	quorum := (total >> 1) + 1
	hashs := make(map[string][]string)

	for _, info := range nodes {
		go func(info cluster.NodeInfo) {
			client := NewNodeClient(info)
			value, err := client.GetHash(key)
			if err != nil {
				result <- hash_node{
					hash: nil,
					node: info.NodeId,
				}
			} else {
				result <- hash_node{
					hash: value,
					node: info.NodeId,
				}
			}
		}(info)
	}

	counter := 0
	for {
		select {
		case hash_node := <-result:

			if hash_node.hash != nil {
				h := fmt.Sprintf("%x", hash_node.hash)
				hashs[h] = append(hashs[h], hash_node.node)
				if len(hashs[h]) >= quorum {
					idx := time.Now().Second() % len(hashs[h])
					client := NewNodeClient(nodes[hashs[h][idx]])
					value, err := client.Get(key)
					if err != nil {
						c.logger.Error("client: get key %s with error %s", key, err)
						return nil, err
					}
					return c.compressor.Decompress(value)
				}
			}

			counter++
			if counter >= total {
				c.logger.Trace("client: key %s is not found", key)
				return nil, ErrNoQuorum
			}

		case <-time.After(c.options.timeout):
			c.logger.Error("client: key %s is got with timeout", key)
			return nil, ErrTimeout
		}
	}

}

func (c *Client) Del(key string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	nodes, err := c.subscriber.GetNodes()
	if err != nil {
		return false, err
	}

	result := make(chan bool)

	total := len(nodes)
	quorum := (total >> 1) + 1
	ok, bad := 0, 0

	for _, info := range nodes {
		go func(info cluster.NodeInfo) {
			client := NewNodeClient(info)
			if status, err := client.Del(key); err == nil && status {
				result <- true
			} else {
				result <- false
			}
		}(info)
	}

	for {
		select {
		case status := <-result:
			if status {
				ok++
			} else {
				bad++
			}
			if ok >= quorum {
				c.logger.Trace("client: key %s is deleted with success", key)
				return true, nil
			}
			if bad >= quorum {
				c.logger.Trace("client: key %s is deleted with error", key)
				return false, ErrNoQuorum
			}
			if ok+bad >= total {
				c.logger.Trace("client: key %s is deleted with error", key)
				return false, ErrNoQuorum
			}
		case <-time.After(c.options.timeout):
			c.logger.Error("client: key %s is deleted with timeout", key)
			return false, ErrTimeout
		}
	}
}
