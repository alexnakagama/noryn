package project

import (
	"os"
	"path/filepath"
)

type Context struct {
	Root  string
	Files []string
}

func (p *Project) BuildContext() (Context, error) {
	var files []string

	err := filepath.Walk(p.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(p.Root, path)
		if err != nil {
			return err
		}

		files = append(files, relPath)

		return nil
	})

	if err != nil {
		return Context{}, err
	}

	return Context{
		Root:  p.Root,
		Files: files,
	}, nil
}
