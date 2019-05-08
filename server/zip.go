package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateZipOf(path string) {
	log.Printf("Creating zip of %s\n", config.Root + path)

	// Pick the spot for the new zip
	tmp := config.Root + path + ".zip"
	zipFile, err := os.Create(tmp)
	if err != nil { log.Fatal(err) }
	defer zipFile.Close()

	// Initialize the writer
	zipFileWriter := zip.NewWriter(zipFile)
	defer zipFileWriter.Close()

	filepath.Walk(config.Root + path, addToZip(zipFileWriter))
}

func addToZip(zipFile *zip.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }

		// Including directories screws up the unzipping process
		if info.IsDir() { return nil }

		reader, err := os.Open(path)

		header, err := zip.FileInfoHeader(info)
		if err != nil { log.Fatal(err) }

		header.Name = strings.Replace(path, config.Root + "/", "", 1)
		header.Method = zip.Deflate

		writer, err := zipFile.CreateHeader(header)
		if err != nil { log.Fatal(err) }

		_, err = io.Copy(writer, reader)
		return nil
	}
}
