package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

type skipPredicate func(d fs.DirEntry) (skip bool)

// walkDir walks through a specified directory and sends file paths into the channel of results.
func walkDir(dir string, skip skipPredicate) (result chan string, status chan error) {
	result = make(chan string)
	status = make(chan error)

	go func() {
		status <- filepath.WalkDir(dir, func(p string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}

			if p == dir || skip(d) {
				return nil
			}

			if d.IsDir() {
				return fs.SkipDir
			}
			result <- p
			return nil
		})
		close(result)
		close(status)
	}()

	return result, status
}

func transformVal(val []byte) string {
	return strings.TrimRight(
		string(
			bytes.ReplaceAll(val, []byte{0x00}, []byte("\n")),
		),
		" \t",
	)
}

func readValueFromFile(path string) (EnvValue, error) {
	stat, _ := os.Stat(path)
	if stat.Size() == 0 {
		return EnvValue{Value: "", NeedRemove: true}, nil
	}

	f, e := os.Open(path)
	if e != nil {
		return EnvValue{}, e
	}
	defer f.Close()

	lineBytes, _, e := bufio.NewReader(f).ReadLine()
	if e != nil && !errors.Is(e, io.EOF) {
		return EnvValue{}, e
	}

	return EnvValue{Value: transformVal(lineBytes)}, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (res Environment, err error) {
	res = make(Environment)
	files, status := walkDir(dir, func(d fs.DirEntry) (skip bool) {
		return strings.Contains(d.Name(), "=")
	})

	for {
		select {
		case file := <-files:
			line, e := readValueFromFile(file)
			if e != nil {
				continue
			}
			res[filepath.Base(file)] = line
		case err = <-status:
			return res, err
		}
	}
}
