package project

import (
	"errors"
	"os"
	"path/filepath"
)

func FindRoot(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)

		if parent == current {
			return "", errors.New("project root not found")
		}

		current = parent
	}
}

func Discover(start string) (*Project, error) {
	root, err := FindRoot(start)
	if err != nil {
		return nil, err
	}

	return &Project{
		Root: root,
	}, nil
}
