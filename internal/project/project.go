package project

import (
	"fmt"
	"path/filepath"
)

type Project struct {
	Root string
}

func (p *Project) ResolvePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	absPath, err := filepath.Abs(filepath.Join(p.Root, path))
	if err != nil {
		return "", err
	}

	relPath, err := filepath.Rel(p.Root, absPath)
	if err != nil {
		return "", err
	}

	if relPath == ".." || len(relPath) >= 3 && relPath[:3] == ".."+string(filepath.Separator) {
		return "", fmt.Errorf("path is outside the project: %s", path)
	}

	return absPath, nil
}
