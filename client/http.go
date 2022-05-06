package client

import (
	"net"
	"net/http"
	"time"
)

const (
	// DefaultTimeout 默认的HTTP请求超时
	DefaultTimeout = 3000 // 3s
	// DefaultKeepAlive 默认的KeepAlive时间
	DefaultKeepAlive = 60000 // 60s
	// DefaultMaxIdleConns 默认的最大空闲连接数
	DefaultMaxIdleConns = 1000
	// DefaultMaxIdleConnsPerHost 默认的与每台机器的最大空闲连接数
	DefaultMaxIdleConnsPerHost = 100
)

// HTTPConfig HTTP配置
type HTTPConfig struct {
	Timeout             int `json:"timeout" mapstructure:"timeout"`
	KeepAlive           int `json:"keepAlive" mapstructure:"keepAlive"`
	MaxIdleConns        int `json:"maxIdleConns" mapstructure:"maxIdleConns`
	MaxIdleConnsPerHost int `json:"maxIdleConnsPerHost" mapstructure:"maxIdleConnsPerHost"`
}

// NewHTTPClient 返回一个原生的http.Client
func NewHTTPClient(conf *HTTPConfig) *http.Client {
	timeout := conf.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	keepAlive := conf.KeepAlive
	if keepAlive <= 0 {
		keepAlive = DefaultKeepAlive
	}

	maxIdleConns := conf.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = DefaultMaxIdleConns
	}
	maxIdleConnsPerHost := conf.MaxIdleConnsPerHost
	if maxIdleConnsPerHost <= 0 {
		maxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   time.Duration(timeout) * time.Millisecond,
				KeepAlive: time.Duration(keepAlive) * time.Millisecond,
			}).Dial,
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
		},
	}
}
