package tui

type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenChat
)

type Router struct {
	current Screen
}

func NewRouter() Router {
	return Router{
		current: ScreenWelcome,
	}
}

func (r *Router) Current() Screen {
	return r.current
}

func (r *Router) Navigate(screen Screen) {
	r.current = screen
}
