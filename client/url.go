package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func Format(val interface{}) (string, error) {
	var err error
	switch v := val.(type) {
	case string:
		return v, err
	case int:
		return strconv.Itoa(v), err
	case int32:
		return strconv.FormatInt(int64(v), 10), err
	case int64:
		return strconv.FormatInt(v, 10), err
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), err
	case bool:
		return strconv.FormatBool(v), err
	default:
		return fmt.Sprintf("%v", v), err
	}
}

func Convert(params interface{}) (url.Values, error) {
	if values, ok := params.(url.Values); ok {
		return values, nil
	}
	content, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	var maps map[string]interface{}
	if err = json.Unmarshal(content, &maps); err != nil {
		return nil, err
	}
	values := url.Values{}
	for key, val := range maps {
		var value string
		switch t := val.(type) {
		case int, int32, int64, float32, float64, string, bool:
			str, err := Format(t)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %T(%v=%v) to string: %v", val, key, val, err)
			}
			value = str
		default:
			return nil, fmt.Errorf("unexpected type %T while converting: %v=%v", t, key, val)
		}
		values.Add(key, value)
	}
	return values, nil
}
