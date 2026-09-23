package project

import (
	"fmt"
	"os"
	"path/filepath"
)

type Project struct {
	Root string
}

func (p *Project) ResolvePath(path string) (string, error) {
	fmt.Printf("ROOT: %q\n", p.Root)
	fmt.Printf("PATH: %q\n", path)

	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	root, err := filepath.Abs(p.Root)
	if err != nil {
		return "", err
	}

	var absPath string

	if filepath.IsAbs(path) {
		absPath = filepath.Clean(path)
	} else {
		absPath, err = filepath.Abs(filepath.Join(root, path))
		if err != nil {
			return "", err
		}
	}

	relPath, err := filepath.Rel(root, absPath)
	if err != nil {
		return "", err
	}

	if relPath == ".." ||
		len(relPath) >= 3 &&
			relPath[:3] == ".."+string(filepath.Separator) {
		return "", fmt.Errorf("path is outside the project: %s", path)
	}

	if _, err := os.Stat(absPath); err != nil {
		return "", err
	}

	return absPath, nil
}
