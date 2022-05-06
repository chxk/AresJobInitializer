package utils

import (
	"encoding/json"
	"github.com/golang/glog"
)

func Unmarshal(content []byte, object interface{}) error {
	if object == nil {
		return nil
	}
	if err := json.Unmarshal(content, object); err != nil {
		glog.Errorf("failed to unmarshal response body to %T: body=%s, error=%v", object, content[:50], err)
		return err
	}
	return nil
}
