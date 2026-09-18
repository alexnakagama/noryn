package instructions

import (
	"os"

	"github.com/alexnakagama/noryn/internal/project"
)

type Builder struct {
	project *project.Project
}

func NewBuilder(project *project.Project) *Builder {
	return &Builder{
		project: project,
	}
}

func (b *Builder) Build() (string, error) {
	path, err := b.project.ResolvePath("AGENTS.md")
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
