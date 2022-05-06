package client

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	values := url.Values{}
	values.Add("id", "1000000000")
	values.Add("id", "1000000001")

	vals, err := Convert(values)
	assert.Nil(t, err)
	assert.Equal(t, values, vals)
	t.Logf("vals: %#v", vals)

	var float32Val float32 = 0.0000001
	maps := map[string]interface{}{
		"float": float32Val,
	}
	values = url.Values{}
	values.Add("float", "0.0000001")
	vals, err = Convert(maps)
	assert.Nil(t, err)
	assert.Equal(t, values, vals)
	t.Logf("vals: %#v", vals)
}
