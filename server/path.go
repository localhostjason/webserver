package server

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetExeDir() (string, error) {
	dir, err := getCurrentAbPathByExecutable()
	if err != nil {
		return "", err
	}

	if strings.Contains(dir, getTmpDir()) {
		return getCurrentAbPathByCaller(), nil
	}
	return dir, nil
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	res, err := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res, err
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
