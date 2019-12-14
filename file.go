package lade

import (
	"io"
	"os"
	"path/filepath"

	"github.com/containers/storage/pkg/archive"
	"github.com/zealic/xignore"
)

type File struct {
	Body io.Reader
	Name string
}

func GetTarFile() (*File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	exclude := new(xignore.Ignorefile)
	input, err := os.Open(filepath.Join(cwd, ".dockerignore"))
	if err == nil {
		err = exclude.FromReader(input)
		if err != nil {
			return nil, err
		}
	}
	exclude.Patterns = append(exclude.Patterns, ".git", ".bzr", ".hg", ".svn")
	opts := &archive.TarOptions{
		Compression:     archive.Zstd,
		ExcludePatterns: exclude.Patterns,
		IncludeFiles:    []string{"."},
	}
	body, err := archive.TarWithOptions(cwd, opts)
	if err != nil {
		return nil, err
	}
	file := new(File)
	file.Body = body
	file.Name = filepath.Base(cwd) + "." + opts.Compression.Extension()
	return file, nil
}
