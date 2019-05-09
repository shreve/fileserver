package main

import (
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
)

type File struct {
	Name string   `json:"name"`
	Dir bool      `json:"dir"`
	Zipped bool   `json:"zipped"`
	Size int64    `json:"size"`
}

func GetAllFiles(path string) ([]File, error) {
	entries, _ := ioutil.ReadDir(config.Root + path)

	files := []File{}
	for _, info := range entries {
		file := File{
			info.Name(),
			info.IsDir(),
			zipExists(path, info),
			size(config.Root + path, info)}
		files = append(files, file)
	}
	return files, nil
}

func zipExists(path string, info os.FileInfo) bool {
	if (!info.IsDir()) { return false }
	zip := config.TmpDir + path + info.Name() + ".zip"
	log.Printf("Zip to test: %s", zip)
	if _, err := os.Stat(zip); err != nil {
		return !os.IsNotExist(err)
	}
	return false
}

func size(path string, info os.FileInfo) int64 {
	if (!info.IsDir()) { return info.Size() }
	var out int64 = 0
	filepath.Walk(path + info.Name(), func(path string, file os.FileInfo, err error) error {
		if (err != nil) { return err }
		out += size(path, file)
		return nil
	})
	return out
}
