package tools

type Tool interface {
	Name() string
	Execute(arguments string) (string, error)
}
