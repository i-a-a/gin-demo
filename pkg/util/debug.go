package util

import (
	"fmt"
	"runtime"
	"strings"
)

// 获取错误文件和行号。 去除go自带函数和外部包。
func GetErrorStack(errString string, splitDirName string) []string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	idx := 0
	recorder := []string{}

	for i, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		filepath, line := fn.FileLine(pc)

		if strings.Contains(filepath, splitDirName) {
			if idx == 0 {
				idx = strings.Index(filepath, splitDirName)
			}
			recorder = append(recorder, fmt.Sprintf("%s:%d", filepath[idx:], line))
		}

		if i >= 20 {
			break
		}
	}

	return recorder
}
