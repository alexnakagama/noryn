package project

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var ignoredDirectories = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
}

type Context struct {
	Root  string
	Files []string
}

func shouldIncludeFile(path string, info os.FileInfo) bool {}

func (p *Project) BuildContext() (Context, error) {
	var files []string

	err := filepath.Walk(p.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if ignoredDirectories[info.Name()] {
				return filepath.SkipDir
			}

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

	sort.Strings(files)

	return Context{
		Root:  p.Root,
		Files: files,
	}, nil
}

func (c Context) String() string {
	var builder strings.Builder

	builder.WriteString("Project root: ")
	builder.WriteString(c.Root)
	builder.WriteString("\nFiles:\n")

	for _, file := range c.Files {
		builder.WriteString("- ")
		builder.WriteString(file)
		builder.WriteString("\n")
	}

	return builder.String()
}
