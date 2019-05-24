package main

import (
	"log"
	"fmt"
	"os"
	"time"
	"io/ioutil"
	"path/filepath"
)

type File struct {
	Name string   `json:"name"`
	Dir bool      `json:"dir"`
	Zipped bool   `json:"zipped"`
	Size int64    `json:"size"`
}

type Unit int

const (
	Bytes Unit = iota
	Kilobytes = iota
	Megabytes = iota
	Gigabytes = iota
	Terabytes = iota
)

type BitRate struct {
	bytes int64
	period time.Duration
}

func (b *BitRate) Human() string {
	bps := int64(float64(b.bytes) / b.period.Seconds())
	dur := time.Duration(1) * time.Second
	avg := BitRate{bps, dur}
	return fmt.Sprintf("%sps", avg.Size())
}

func (b *BitRate) Size() string {
	size, unit := b.humanSize()
	label := ""
	switch unit {
	case Bytes: label = "B"
	case Kilobytes: label = "KB"
	case Megabytes: label = "MB"
	case Gigabytes: label = "GB"
	case Terabytes: label = "TB"
	}
	return fmt.Sprintf("%.2f%s", size, label)
}

func (b *BitRate) humanSize() (float64, Unit) {
	val := float64(b.bytes)
	unit := Bytes
	for val > 1024 {
		val /= 1024
		unit++
	}
	return val, Unit(unit)
}

func GetAllFiles(path string) ([]File, error) {
	entries, _ := ioutil.ReadDir(config.Root + path)

	files := []File{}
	for _, info := range entries {
		if info.Mode() & os.ModeSymlink != 0 {
			continue;
		}
		file := File{
			info.Name(),
			info.IsDir(),
			zipExists(path, info),
			Size(config.Root + path, info)}
		files = append(files, file)
	}
	return files, nil
}

func DownloadFile(path string) (string, error) {
	log.Printf("Preparing to download %s", path)

	info, err := os.Stat(config.Root + path)
	if err != nil { return "", err }

	if info.IsDir() {
		err = SafeCreateZipOf(path)
		if err != nil { return "", err }
		return config.TmpDir + path + ".zip", nil
	}

	return config.Root + path, nil
}

func StatusFile(path string) error {
	log.Printf("Checking status of %s", path)

	if _, err := os.Stat(config.Root + path); err != nil {
		return err
	}

	if err := IsZipping(path); err != nil {
		return err
	}

	if !ZipExists(path) {
		return ZipNotBeingCreated;
	}

	return nil
}

func zipExists(path string, info os.FileInfo) bool {
	if (!info.IsDir()) { return false }
	zip := config.TmpDir + path + info.Name() + ".zip"
	log.Printf("Zip to test: %s", zip)
	if _, err := os.Stat(zip); err == nil {
		return true
	}
	return false
}

func Size(path string, info os.FileInfo) int64 {
	if (!info.IsDir()) { return info.Size() }
	var out int64 = 0
	filepath.Walk(path + "/" + info.Name(), func(path string, file os.FileInfo, err error) error {
		if (err != nil) { return err }
		out += Size(path, file)
		return nil
	})
	return out
}
