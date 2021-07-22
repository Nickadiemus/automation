package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"regexp"
)

// by defalt we're filtering out file ext .mp4, .webn, .DS_Store, .zip
var re = regexp.MustCompile(`^.*\.(MP4|mp4|webm|WEBM|DS_Store|zip|ZIP)$`)

// Handles file compression logic
func ZipFiles(name, relZipDir string, files []os.FileInfo) error {
	// create new archive file name
	zFile, err := os.Create(name)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer zFile.Close() // ensure file closes
	// create new zip writer
	zWriter := zip.NewWriter(zFile)
	// loop over files and add call the zip writer function
	defer zWriter.Close()
	for _, file := range files {
		if !re.MatchString(file.Name()) && !file.IsDir() {
			// TODO: Write concurrrently (but make sure to no race cases )
			if err := WriteFile(relZipDir+file.Name(), zWriter); err != nil {
				return err
			}
		}
	}
	return nil
}

// Writes provided file to zip writer
func WriteFile(file string, w *zip.Writer) error {
	// open current file
	wFile, err := os.Open(file)
	if err != nil {
		return err
	}
	// ensure to close resource
	defer wFile.Close()

	wr, err := w.Create(file)

	if err != nil {
		return err
	}
	_, err = io.Copy(wr, wFile)
	return err
}
