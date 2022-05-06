package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/glog"
)

// Config Client配置
type Config struct {
	Host  string     `json:"host" mapstructure:"host"`
	Token string     `json:"token" mapstructure:"token"`
	HTTP  HTTPConfig `mapstructure:"http"`
}

// Client HTTP客户端
type Client struct {
	Host               string
	ContentType        string
	DisableVersboseLog bool
	Header             http.Header
	HTTPClient         *http.Client
}

// NewClient 返回一个新的客户端
func NewClient(conf *Config) *Client {
	return &Client{
		Host:        conf.Host,
		ContentType: "application/json",
		Header:      make(http.Header),
		HTTPClient:  NewHTTPClient(&conf.HTTP),
	}
}

// Clone 生成一个新的Client，但是会共享同一个HTTP Client
func (cli *Client) Clone() *Client {
	newCli := *cli
	newCli.Header = cli.Header.Clone()
	return &newCli
}

// WithHeader 附加额外的Header到后续的HTTP请求中
func (cli *Client) WithHeader(key, value string) *Client {
	newCli := cli.Clone()
	newCli.Header.Add(key, value)
	return newCli
}

// WithoutVerboseLog 设置不打印详细日志
func (cli *Client) WithoutVerboseLog() *Client {
	newCli := cli.Clone()
	newCli.DisableVersboseLog = true
	return newCli
}

// Do 发送HTTP请求
func (cli *Client) Do(method string, path string, params url.Values, data interface{}) ([]byte, error) {
	// prepare request
	var rb io.Reader
	if data != nil {
		content, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		rb = bytes.NewReader(content)
	}

	url := cli.Host + path
	req, err := http.NewRequest(method, url, rb)
	if err != nil {
		return nil, err
	}
	req.Header = cli.Header.Clone()
	req.Header.Add("Content-Type", cli.ContentType)
	req.URL.RawQuery = params.Encode()

	// do request
	start := time.Now()
	resp, err := cli.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !cli.DisableVersboseLog {
		glog.Infof("\"%s %s\": code=%v, consumed=%v", req.Method, req.URL.String(), resp.StatusCode, time.Now().Sub(start))
	}

	// handle response
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		var errorInfo HTTPError
		if err := json.Unmarshal(content, &errorInfo); err == nil {
			errorInfo.HTTPCode = resp.StatusCode
			if errorInfo.Code == "" && errorInfo.Status != "" {
				errorInfo.Code = errorInfo.Status
			}
			return nil, &errorInfo
		}
		return nil, fmt.Errorf("[%d] %s", resp.StatusCode, content[:50])
	}
	return content, nil
}

// DoRequest 发送HTTP请求
func (cli *Client) DoRequest(method string, path string, params interface{}, data interface{}) ([]byte, error) {
	if params == nil {
		return cli.Do(method, path, nil, data)
	}
	if values, ok := params.(url.Values); ok {
		return cli.Do(method, path, values, data)
	}
	values, err := Convert(params)
	if err != nil {
		return nil, fmt.Errorf("failed to convert query params to url.Values: %v", err)
	}
	return cli.Do(method, path, values, data)
}

/*
examples #1:
	cli.Get("/url", 0, nil)

examples #2:
	params := map[string]string{"key": "value"}
	cli.Get("/url", 0, params)

examples #3:
	type Query struct {
		Key string `json:"key,omitempty"`
	}
	params := Query{Key: "value"}
	cli.Get("/url", 0, params)
*/
// Get 发送GET请求
func (cli *Client) Get(path string, params interface{}) ([]byte, error) {
	return cli.DoRequest(http.MethodGet, path, params, nil)
}

// Post 发送POST请求
func (cli *Client) Post(path string, params interface{}, data interface{}) ([]byte, error) {
	return cli.DoRequest(http.MethodPost, path, params, data)
}

// Patch 发送PATCH请求
func (cli *Client) Patch(path string, params interface{}, data interface{}) ([]byte, error) {
	return cli.DoRequest(http.MethodPatch, path, params, data)
}

// Put 发送PUT请求
func (cli *Client) Put(path string, params interface{}, data interface{}) ([]byte, error) {
	return cli.DoRequest(http.MethodPut, path, params, data)
}

// Delete 发送Delete请求
func (cli *Client) Delete(path string, params interface{}) ([]byte, error) {
	return cli.DoRequest(http.MethodDelete, path, params, nil)
}
