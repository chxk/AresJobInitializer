package cache

import (
	"net/url"
	"os"
	"testing"

	"opensource/aresjob-initializer/utils"

	"github.com/stretchr/testify/assert"
)

func TestURI(t *testing.T) {
	uri := "nodePort://127.0.0.1,127.0.0.2:30037"
	u, err := url.Parse(uri)
	assert.Nil(t, err)
	t.Logf("u: %#v", u)
}

func TestParse(t *testing.T) {
	os.Setenv(utils.LOCAL_IP, "0.0.0.0")
	uri := "127.0.0.1,127.0.0.2:30037"
	ips, port, err := parseNodePortURI(uri)
	assert.Nil(t, err)
	assert.Equal(t, []string{"0.0.0.0", "127.0.0.1", "127.0.0.2"}, ips)
	assert.Equal(t, "30037", port)
	t.Logf("ips: %v, port: %v", ips, port)

	uri = "30037"
	ips, port, err = parseNodePortURI(uri)
	assert.Nil(t, err)
	assert.Equal(t, []string{"0.0.0.0"}, ips)
	assert.Equal(t, "30037", port)
	t.Logf("ips: %v, port: %v", ips, port)

	uri = ""
	_, _, err = parseNodePortURI(uri)
	assert.NotNil(t, err)
	t.Logf("err: %v", err)
}
