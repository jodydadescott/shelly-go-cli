package types

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	STDIN    bool
	BaseName string
	FullName string
	Bytes    []byte
}

type Files struct {
	IsDir bool
	File  *File
	Files []*File
}

func newFile(name string, bytes []byte, stdin bool) *File {
	f := &File{
		STDIN:    stdin,
		FullName: name,
		Bytes:    bytes,
	}

	if !stdin {
		f.BaseName = filepath.Base(name)
	}

	return f
}

func (t *Files) GetFiles(name string) *File {

	for _, file := range t.Files {
		if !file.STDIN {
			if file.BaseName == name {
				return file
			}
			if file.BaseName == name+".yaml" {
				return file
			}
			if file.BaseName == name+".json" {
				return file
			}

			split := strings.Split(name, "-")

			if file.BaseName == split[0]+".yaml" {
				return file
			}

			if file.BaseName == split[0]+".json" {
				return file
			}

		}
	}

	for _, file := range t.Files {
		if file.STDIN {
			return file
		}
	}

	return nil
}

func NewFiles(input string) (*Files, error) {

	files := &Files{}

	if input != "" {

		info, err := os.Stat(input)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {

			files.IsDir = true

			rawFiles, err := ioutil.ReadDir(input)
			if err != nil {
				return nil, err
			}

			for _, rawFile := range rawFiles {
				bytes, err := os.ReadFile(input + "/" + rawFile.Name())
				if err != nil {
					return nil, err
				}

				files.Files = append(files.Files, newFile(rawFile.Name(), bytes, false))
			}

		} else {

			bytes, err := os.ReadFile(input)
			if err != nil {
				return nil, err
			}

			files.Files = append(files.Files, newFile(input, bytes, false))

		}

		return files, nil
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {

		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}

		files.Files = append(files.Files, newFile("", bytes, true))

		return files, nil
	}

	return nil, fmt.Errorf("Data input is required. Use filename or pipe to STDIN")
}
