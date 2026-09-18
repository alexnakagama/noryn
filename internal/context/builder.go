package context

import "github.com/alexnakagama/noryn/internal/project"

type Builder struct {
	project *project.Project
}

func NewBuilder(project *project.Project) *Builder {
	return &Builder{
		project: project,
	}
}

func (b *Builder) Build() (string, error) {}
