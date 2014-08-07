package main

type Handler struct {
	CommandName    string
	CommandPattern string
	CommandParser  func(*Command) (map[string]string, bool)
	HandlerFunc    func(*Command) string
}

func NewHandler() *Handler {
	return &Handler{}
}
