package util

import (
	"os"
	"path/filepath"
)

func OpenFilehandle(path string) (*os.File, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	root := filepath.Dir(abs_path)

	_, err = os.Stat(root)

	if os.IsNotExist(err) {

		err := os.MkdirAll(root, 0755)

		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
}
