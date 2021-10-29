package utils

import (
	"os"
	"path/filepath"
)

const (
	ConfigEnv  = "V6_CONFIG"
	ConfigFile = "config.yaml"
)

func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()

}

func GetAllFile(pathname string) (files []string) {
	root := pathname
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
