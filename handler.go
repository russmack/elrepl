package main

type Handler struct {
	CommandName    string
	CommandPattern string
	Usage          string
	CommandParser  func(*Command) (map[string]string, bool)
	CommandParse   func(*Command) (ParseMap, bool)
	HandlerFunc    func(*Command) string
}

func NewHandler() *Handler {
	return &Handler{}
}
