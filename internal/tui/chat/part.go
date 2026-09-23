package chat

type Part interface {
	part()
}

type TextPart struct {
	Content string
}

func (TextPart) part() {}

type ToolPart struct {
	Name      string
	Arguments string
	Status    string
	Result    string
}

func (ToolPart) part() {}

type ErrorPart struct {
	Content string
}

func (ErrorPart) part() {}
