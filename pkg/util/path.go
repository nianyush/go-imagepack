package util

import (
	"os"
	"path/filepath"
)

func AbsPath(path string) (string,error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path = filepath.Join(wd, path)

	return path, nil
}