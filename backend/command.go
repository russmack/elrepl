package backend

import (
	"regexp"
)

type Command struct {
	Instruction string
	Body        string
	Name        string
	Pattern     *regexp.Regexp
	Tokens      []string
	//Args        []string
}

type Resource struct {
	Scheme   string
	Host     string
	Port     string
	Index    string
	Type     string
	Field    string
	Endpoint string
	Alias    string
	Routing  string
}

func NewCommand(instr string, body string, name string, tokens []string) *Command {
	return &Command{Instruction: instr, Body: body, Name: name, Tokens: tokens}
}
