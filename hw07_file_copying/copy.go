package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const (
	BUFFERSIZE = 1024
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return errors.New("fromPath equal toPath")
	}

	if offset < 0 {
		return errors.New("offset < 0")
	}

	if limit < 0 {
		return errors.New("limit < 0")
	}

	fi, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	size := fi.Size()

	ofile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}
	defer ofile.Close()

	wfile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer wfile.Close()

	_, err = ofile.Seek(offset, io.SeekCurrent)
	if err != nil {
		return err
	}
	an := int64(0)
	var reader io.Reader
	if limit > 0 {
		reader = &io.LimitedReader{R: ofile, N: limit}
		size = limit
	} else {
		reader = ofile
	}

	for {
		n, err := io.Copy(wfile, reader)

		an += n
		fmt.Printf("Progress: %3.2f%%", 100*float32(size)/float32(an))

		if !errors.Is(err, nil) && !errors.Is(err, io.EOF) {
			return err
		}
		if n == 0 {
			break
		}
	}

	return nil
}
