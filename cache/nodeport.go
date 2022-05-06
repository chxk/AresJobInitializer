package cache

import (
	"fmt"
	"net/url"
	"strings"

	"aresjob-initializer/utils"

	"github.com/golang/glog"
)

type NodePortClient struct {
	clients []*ObserverClient
}

func NewNodePortClient(schema, host string, config *Config) (*NodePortClient, error) {
	ips, port, err := parseNodePortURI(host)
	if err != nil {
		return nil, err
	}
	client := &NodePortClient{}
	for _, ip := range ips {
		u := url.URL{Scheme: schema, Host: fmt.Sprintf("%s:%s", ip, port)}
		c, err := NewObserverClient(u.String(), config)
		if err != nil {
			return nil, err
		}
		client.clients = append(client.clients, c)
	}
	return client, nil
}

func (c *NodePortClient) GetJobCache(key string) (cache JobCache, err error) {
	for _, cli := range c.clients {
		cache, err = cli.GetJobCache(key)
		if err == nil {
			return cache, nil
		}
		glog.Warningf("%s failed: %v", cli.Host, err)
	}
	return cache, err
}

func (c *NodePortClient) Close() {
	for _, cli := range c.clients {
		cli.Close()
	}
}

func parseNodePortURI(uri string) ([]string, string, error) {
	parts := strings.Split(uri, ":")
	if len(uri) == 0 || len(parts) < 1 {
		return nil, "", fmt.Errorf("unknown nodePort schema: %v", uri)
	}
	localIP := utils.GetLocalIP()
	port := parts[len(parts)-1]
	if len(parts) == 1 {
		return []string{localIP}, port, nil
	}
	parts = strings.Split(parts[0], ",")
	return append([]string{localIP}, parts...), port, nil
}
