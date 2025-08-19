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

func Copy(fromPath, toPath string, offset, limit int64) error {
	var buf []byte
	if limit > 0 {
		buf = make([]byte, limit)
	} else {
		buf = make([]byte, 1024)
	}

	fileFrom, openingError := os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if openingError != nil {
		return openingError
	}
	stat, err := os.Stat(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	size := stat.Size()
	if size < offset {
		return ErrOffsetExceedsFileSize
	}
	copiedSize := size - offset
	if limit < copiedSize && limit != 0 {
		copiedSize = limit
	}
	fileTo, _ := os.Create(toPath)
	currentSize := 0
	for offset < copiedSize {
		n, readingError := fileFrom.Read(buf)
		offset += int64(n)
		if errors.Is(readingError, io.EOF) {
			break
		} else if readingError != nil && !errors.Is(readingError, io.EOF) {
			return readingError
		}

		write, writingError := fileTo.Write(buf[:n])
		fmt.Printf("Current percentage %d%%", int64(currentSize+write)/copiedSize*100)
		currentSize += write
		if writingError != nil {
			return writingError
		}
	}
	defer func() {
		fileFrom.Close()
		fileTo.Close()
	}()

	return nil
}
