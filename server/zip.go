package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"errors"
	"time"
)

var ZipNotBeingCreated = errors.New("The requested file has not been queued for compression");
var ZipBeingCreated = errors.New("The requested file is being compressed for download");

func SafeCreateZipOf(path string) error {
	// Pick the spot for the new zip
	flag := config.TmpDir + path + ".zipping"

	if err := IsZipping(path); err != nil {
		// Detect if the output file is being written to
		return err
	} else if !ZipExists(path) {
		// Otherwise indicate we're starting
		os.Create(flag)
		go func() {
			start := time.Now()
			CreateZipOf(path)

			// Indicate we are done compressing
			os.Remove(flag)

			info, err := os.Stat(config.Root + path)
			if err != nil { log.Fatal(err) }

			diff := time.Now().Sub(start)
			rate := BitRate{Size(config.Root, info), diff}
			log.Printf("Completed compression of %s (%s) in %s at a rate of %s",
				path, rate.Size(), diff, rate.Human())
		}()
		return ZipBeingCreated
	}

	return nil
}

func CreateZipOf(path string) {
	log.Printf("Creating zip of %s\n", config.Root + path)

	// Create a file to write to
	zipPath := config.TmpDir + path + ".zip"
	zipFile, err := os.Create(zipPath)
	if err != nil {
		log.Printf("Making directories down to %s", filepath.Dir(zipPath))
		os.MkdirAll(filepath.Dir(zipPath), os.ModePerm)
		zipFile, err = os.Create(zipPath)
	}
	if err != nil { log.Fatal(err) }
	defer zipFile.Close()

	// Initialize the writer
	zipFileWriter := zip.NewWriter(zipFile)
	defer zipFileWriter.Close()

	// Recursively add every file to the archive
	err = filepath.Walk(config.Root + path,
		addToZip(filepath.Base(path), config.Root + path, zipFileWriter))

	if err != nil {
		log.Printf("Finished with an errror: %s", err)
		os.Remove(zipPath)
	} else {
		log.Printf("Finished creating %s successfully", zipPath)
	}
}

func ZipExists(path string) bool {
	tmp := config.TmpDir + path + ".zip"
	_, err := os.Stat(tmp)
	return err == nil
}

func IsZipping(path string) error {
	flag := config.TmpDir + path + ".zipping"
	if _, err := os.Stat(flag); err == nil {
		return ZipBeingCreated
	}
	return nil
}

func addToZip(foldername string, startpath string, zipFile *zip.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }

		// Including directories screws up the unzipping process
		if info.IsDir() { return nil }

		reader, err := os.Open(path)

		header, err := zip.FileInfoHeader(info)
		if err != nil { log.Fatal(err) }

		header.Name = strings.Replace(path, startpath, foldername, 1)
		header.Method = zip.Deflate

		writer, err := zipFile.CreateHeader(header)
		if err != nil { log.Fatal(err) }

		_, err = io.Copy(writer, reader)
		return nil
	}
}
