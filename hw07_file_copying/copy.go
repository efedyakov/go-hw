package main

import (
	"errors"
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
	checklimit := limit != 0
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

	buf := make([]byte, BUFFERSIZE)
	_, err = ofile.Seek(offset, io.SeekCurrent)
	if err != nil {
		return err
	}
	an := int64(0)

	for {
		n, err := ofile.Read(buf)
		if an+int64(n) > limit && checklimit {
			n = int(limit - an)
		}
		an += int64(n)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := wfile.Write(buf[:n]); err != nil {
			return err
		}

		if n == int(limit) && checklimit {
			break
		}
	}

	return nil
}
