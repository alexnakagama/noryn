package project

import (
	"net/http"
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

func shouldIncludeFile(path string, info os.FileInfo) bool {
	if info.Size() > 1_000_000 {
		return false
	}

	return true
}

func isBinaryFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil {
		return false, err
	}

	contentType := http.DetectContentType(buffer[:n])

	return !strings.HasPrefix(contentType, "text/") &&
		contentType != "application/json" &&
		contentType != "application/xml" &&
		contentType != "application/javascript", nil
}

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

		if !shouldIncludeFile(path, info) {
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
