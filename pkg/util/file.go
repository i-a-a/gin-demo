package util

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	rootDir string
)

func init() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	rootDir, _ = filepath.EvalSymlinks(filepath.Dir(exePath))
	rootDir = strings.TrimSuffix(rootDir, "/tmp") // air
}

// 执行文件所在目录 (go run 不准)
func GetRootDir() string {
	return rootDir
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func CheckDir(dir string) error {
	if !FileExist(dir) {
		if err := os.MkdirAll(dir, 0644); err != nil {
			return errors.New(err.Error() + ": " + dir)
		}
	}
	return nil
}

func MustOpenFile(filename string) *os.File {
	file, err := OpenFile(filename)
	if err != nil {
		panic(err)
	}
	return file
}

func OpenFile(filename string) (*os.File, error) {
	if !FileExist(filename) {
		dir := filename[0:strings.LastIndex(filename, "/")]
		if err := CheckDir(dir); err != nil {
			return nil, err
		}
	}

	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
}
