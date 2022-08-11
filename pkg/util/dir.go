package util

import (
	"os"
	"path/filepath"
)

func DirSize(path string) (int64, error) {
	var dirSize int64
	readSizeFunc := func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !file.IsDir() {
			dirSize += file.Size()
		}
		return nil
	}

	err := filepath.Walk(path, readSizeFunc)
	return dirSize, err
}
