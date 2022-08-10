package pkg

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	Filer   filer
	rootDir string
)

type filer struct{}

// 执行文件所在目录 (go run 不准)
func (filer) GetRootDir() string {
	if rootDir == "" {
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		rootDir, _ = filepath.EvalSymlinks(filepath.Dir(exePath))
		rootDir = strings.TrimSuffix(rootDir, "/tmp") // air
	}
	return rootDir
}

// 文件存在
func (filer) IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 检查目录是否存在，不存在则创建
func (filer) CheckDir(dir string) error {
	if !Filer.IsExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.New(err.Error() + ": " + dir)
		}
	}
	return nil
}

// 打开文件
func (filer) Open(filename string) (*os.File, error) {
	if !Filer.IsExist(filename) {
		dir := filename[0:strings.LastIndex(filename, "/")]
		if err := Filer.CheckDir(dir); err != nil {
			return nil, err
		}
	}

	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
}

// 打开文件
func (filer) MustOpen(filename string) *os.File {
	file, err := Filer.Open(filename)
	if err != nil {
		panic(err)
	}
	return file
}
