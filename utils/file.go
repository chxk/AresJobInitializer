package utils

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// ReadInt: 读指定路径文件中的int值
func ReadInt(path string) (int, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	str := string(content)
	return strconv.Atoi(strings.TrimSpace(str))
}

// Exists: 判断文件是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// CreateFile: 创建文件
func CreateFile(path string) error {
	_, err := os.Create(path)
	return err
}
