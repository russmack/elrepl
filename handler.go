package main

import (
	//"errors"
	//"fmt"
	"regexp"
)

type Handler struct {
	CommandName    string
	CommandPattern string
	Usage          string
	CommandParser  func(*Command) (string, bool)
	HandlerFunc    func(*Command) string
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Tokenize(commandEntry string, patternFnMap map[*regexp.Regexp]func([]string) (string, bool)) (string, bool) {
	for k, v := range patternFnMap {
		mm := k.FindAllStringSubmatch(commandEntry, -1)
		if len(mm) > 0 {
			p, ok := v(mm[0])
			return p, ok
		}
	}
	// Invalid command
	//return errors.New("No pattern matches command.")
	return "", false
}
