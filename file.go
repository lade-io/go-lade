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
		if err = exclude.FromReader(input); err != nil {
			return nil, err
		}
	}
	excludes := []string{".git", ".bzr", ".hg", ".svn"}
	for _, pattern := range exclude.Patterns {
		if pattern != "Dockerfile" {
			excludes = append(excludes, pattern)
		}
	}
	opts := &archive.TarOptions{
		Compression:     archive.Zstd,
		ExcludePatterns: excludes,
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
