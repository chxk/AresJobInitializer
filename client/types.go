package client

import (
	"fmt"
)

// HTTPError 请求错误响应
//   参考API规范文档5.2：https://wiki.corp.kuaishou.com/pages/viewpage.action?pageId=339164087
type HTTPError struct {
	HTTPCode int    `json:"httpCode"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	// Status panda error code
	// panda 使用status 保存error code，api模块无法使用toolbox解析到status这个返回值
	Status string `json:"status"`
	Err    string `json:"error"`
}

func (e *HTTPError) Error() string {
	detail := e.Err
	if len(detail) == 0 {
		detail = e.Status
	}
	return fmt.Sprintf("HTTPCode=%d, %s(%s): %s", e.HTTPCode, e.Message, e.Code, detail)
}
