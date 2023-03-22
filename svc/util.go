package svc

import (
	"errors"
	"os"
)

func panicErr(reason string) {
	panic(errors.New(reason))
}

//func exePath() string {
//	exe, err := os.Executable()
//	if err != nil {
//		panicErr("failed to get exe path " + err.Error())
//	}
//	exePath, err := filepath.Abs(exe)
//	if err != nil {
//		panicErr("failed to get abs exe path " + err.Error())
//	}
//	return exePath
//}

func exePath() string {
	exeDir, err := os.Getwd()
	if err != nil {
		panicErr("failed to get exe path " + err.Error())
	}
	return exeDir
}
