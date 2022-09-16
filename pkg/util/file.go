package util

import (
	"os"
	"strings"
)

// 文件存在
func HasFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 检查目录是否存在，不存在则创建
func CheckDir(dir string) error {
	if HasFile(dir) {
		return nil
	}
	return os.MkdirAll(dir, os.ModePerm)
}

// 追加写文件
func OpenFile(filename string) (*os.File, error) {
	if !HasFile(filename) {
		dir := filename[0:strings.LastIndex(filename, "/")]
		if err := CheckDir(dir); err != nil {
			return nil, err
		}
	}

	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
}
