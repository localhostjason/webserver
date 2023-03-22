package db

import (
	"os"
	"path/filepath"
)

func GetWorkDir() (string, error) {
	// 工作目录 类型 pwd
	exeDir, err := os.Getwd()
	return exeDir, err
}

func GetExeDir() (string, error) {
	// 程序目录， exe 二进制文件路径
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}
