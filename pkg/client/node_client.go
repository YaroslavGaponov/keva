package client

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/YaroslavGaponov/keva/pkg/cluster"
)

type NodeClient struct {
	info   cluster.NodeInfo
	client http.Client
}

func NewNodeClient(info cluster.NodeInfo) *NodeClient {
	return &NodeClient{
		info:   info,
		client: http.Client{},
	}
}

func (c *NodeClient) Set(key string, value []byte) (bool, error) {
	url := fmt.Sprintf("%s://%s:%d/%s/%s", c.info.Schema, c.info.Host, c.info.Port, c.info.Path, key)
	resp, err := c.client.Post(url, "application/x-binary", bytes.NewReader(value))
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}
	return true, nil

}

func (c *NodeClient) GetHash(key string) ([]byte, error) {
	url := fmt.Sprintf("%s://%s:%d/%s/%s/hash", c.info.Schema, c.info.Host, c.info.Port, c.info.Path, key)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	hash, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (c *NodeClient) Get(key string) ([]byte, error) {
	url := fmt.Sprintf("%s://%s:%d/%s/%s", c.info.Schema, c.info.Host, c.info.Port, c.info.Path, key)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *NodeClient) Del(key string) (bool, error) {
	url := fmt.Sprintf("%s://%s:%d/%s/%s", c.info.Schema, c.info.Host, c.info.Port, c.info.Path, key)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return false, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}
	return true, nil
}
