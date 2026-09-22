package tools

import "fmt"

type Registry struct {
	tools map[string]Tool
}

func NewRegistry(toolList ...Tool) *Registry {
	registry := &Registry{
		tools: make(map[string]Tool),
	}

	for _, tool := range toolList {
		registry.Register(tool)
	}

	return registry
}

func (r *Registry) Register(tool Tool) error {
	name := tool.Name()

	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool already registered: %s", name)
	}

	r.tools[name] = tool

	return nil
}

func (r *Registry) Get(name string) (Tool, bool) {
	tool, ok := r.tools[name]

	return tool, ok
}
