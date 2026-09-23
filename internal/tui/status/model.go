package status

type Model struct {
	model  string
	status string
	tokens int
}

func New() Model {
	return Model{
		status: "ready",
	}
}

func (m Model) SetStatus(status string) Model {
	m.status = status

	return m
}

func (m Model) Status() string {
	return m.status
}
