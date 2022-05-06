package client

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}
	client := NewClient(config)

	path := "/s"
	params := map[string]string{"wd": "query"}
	// GET http://www.baidu.com/s?wd=query
	content, err := client.Get(path, params)
	assert.Nil(t, err)
	t.Logf("response: %v", string(content)[:15])

	values := url.Values{}
	values.Add("wd", "query")
	content, err = client.Get(path, values)
	assert.Nil(t, err)
}

func TestClientWith(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}
	client := NewClient(config)
	client1 := client.WithHeader("abc", "1")
	client2 := client.WithHeader("def", "2")
	client3 := client2.WithHeader("wwww", "3")
	client3 = client3.WithHeader("wwww", "4")
	fmt.Printf("%#v", client3.Header)
	assert.Len(t, client.Header, 0)
	assert.Len(t, client1.Header, 1)
	assert.Len(t, client2.Header, 1)
	assert.Len(t, client3.Header, 2)
}

func TestClientDoGetWithBody(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}

	client := NewClient(config)
	values := url.Values{}
	values.Add("wd", "query")

	rb, err := client.Do("GET", "/s", values, "body xxxxx")
	assert.Nil(t, err)
	assert.NotEmpty(t, rb)
	t.Logf("Receive body : %s", string(rb))
}

func TestClientDoPostWithNilBody(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}

	client := NewClient(config)
	values := url.Values{}
	values.Add("wd", "query")

	rb, err := client.Do("POST", "/s", values, nil)
	fmt.Println("content: ", string(rb[:50]))
	assert.Nil(t, err)
	assert.NotEmpty(t, rb)
	t.Logf("Receive body : %s\n", string(rb))
}

func TestClientDoPostWithBody(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}

	client := NewClient(config)
	values := url.Values{}
	values.Add("wd", "query")

	rb, err := client.Do("POST", "/s", values, "test body")
	fmt.Println("content: ", string(rb[:50]))
	assert.Nil(t, err)
	assert.NotEmpty(t, rb)
	t.Logf("Receive body : %s\n", string(rb))
}
