package util

import "os"

func PathExists(path string) bool {
	f, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	return f.IsDir()
}
