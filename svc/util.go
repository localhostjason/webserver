package svc

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

func panicErr(reason string) {
	panic(errors.New(reason))
}

func exePath() string {
	exe, err := os.Executable()
	if err != nil {
		panicErr("failed to get exe path " + err.Error())
	}
	exePath, err := filepath.Abs(exe)
	if err != nil {
		panicErr("failed to get abs exe path " + err.Error())
	}
	return exePath
}

func exePathFile() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		panicErr("failed to get exe file " + err.Error())
	}

	return file
}
