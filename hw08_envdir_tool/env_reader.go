package main

import (
	"bufio"
	"io"
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

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.Contains(file.Name(), "=") {
			continue
		}

		f, err := os.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		r := bufio.NewReader(f)
		data, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
		/*data = strings.TrimRightFunc(data, func(r rune) bool {
			switch r {
			case '\t', '\n', ' ':
				return true
			}
			return false
		})*/
		data = strings.TrimSpace(data)
		data = strings.Map(func(r rune) rune {
			if r == 0x00 {
				return '\n'
			}
			return r
		}, data)

		env[file.Name()] = EnvValue{Value: data, NeedRemove: fi.Size() == 0}
	}

	return env, nil
}
