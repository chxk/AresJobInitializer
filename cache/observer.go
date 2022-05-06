package cache

import (
	"fmt"

	"aresjob-initializer/client"
	"aresjob-initializer/utils"

	"github.com/golang/glog"
)

type ObserverClient struct {
	*client.Client
}

func NewObserverClient(host string, config *Config) (*ObserverClient, error) {
	return &ObserverClient{Client: client.NewClient(&client.Config{
		Host: host,
		/*
			&http.Client{
					Transport: &http.Transport{
						Dial: (&net.Dialer{
							Timeout:   time.Duration(timeout) * time.Millisecond,
							KeepAlive: time.Duration(keepAlive) * time.Millisecond,
						}).Dial,
						MaxIdleConns:        maxIdleConns,
						MaxIdleConnsPerHost: maxIdleConnsPerHost,
					},
				}
		*/
		HTTP: client.HTTPConfig{
			Timeout:             int(config.DialTimeout.Milliseconds()),
			KeepAlive:           3000, // 3s
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
		},
	})}, nil
}

func (c *ObserverClient) GetJobCache(key string) (JobCache, error) {
	path := fmt.Sprintf("/cache/keys/%s", key)
	content, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	glog.Infof("succeeded to get job cache: <%s> %s", key, content)
	jobCache := JobCache{}
	err = utils.Unmarshal(content, &jobCache)
	return jobCache, err
}

func (c *ObserverClient) Close() {
}
