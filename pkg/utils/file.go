package utils

import (
	"io/ioutil"
	"os"
)

// ListDir -
func ListDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}

// IsFileNotExist -
func IsFileNotExist(file string) bool {
	_, err := os.Stat(file)
	return os.IsNotExist(err)
}

// IsFolderNotExist -
func IsFolderNotExist(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return true
	}
	return !info.IsDir()
}

// DeleteFile -
func DeleteFile(file string) error {
	err := os.Remove(file)
	return err
}

// DeleteFolder -
func DeleteFolder(dir string) error {
	err := os.RemoveAll(dir)
	return err
}
