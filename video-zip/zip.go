package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
)

func ZipFiles(name string, files []os.FileInfo) error {
	// create new archive file name
	zFile, err := os.Create(name)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer zFile.Close() // ensure file closes
	// create new zip writer
	zWriter := zip.NewWriter(zFile)
	// loop over files and add each file to the zip writer
	defer zWriter.Close()
	for _, file := range files {
		if err := WriteFile(file.Name(), zWriter); err != nil {
			return err
		}
	}
	return nil
}

// chose a separate function to
func WriteFile(file string, w *zip.Writer) error {
	fmt.Println("opening/writing", file)
	// open current file
	wFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer wFile.Close()

	wr, err := w.Create(file)
	// // fetch file info
	// fileInfo, err := wFile.Stat()
	// if err != nil {
	// 	return err
	// }
	// // Create zip file header
	// header, err := zip.FileInfoHeader(fileInfo)
	// if err != nil {
	// 	return err
	// }
	// // Using FileInfoHeader() above only uses the basename of the file. If we want
	// // to preserve the folder structure we can overwrite this with the full path.
	// header.Name = file

	// // Chosing compression method
	// // Change to deflate to gain better compression
	// // see http://golang.org/pkg/archive/zip/#pkg-constants
	// header.Method = zip.Deflate

	if err != nil {
		return err
	}
	_, err = io.Copy(wr, wFile)
	return err
}
