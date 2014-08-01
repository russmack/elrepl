package main

import (
	"regexp"
	"strings"
)

type CommandParser struct {
	commandRegexpMap map[string]*regexp.Regexp
}

func NewCommandParser() *CommandParser {
	return &CommandParser{}
}

func (p *CommandParser) Parse(entry string) (*Command, error) {
	parts := strings.SplitN(entry, " ", 2)
	cmdName := parts[0]
	cmdArgs := ""
	if len(parts) > 1 {
		cmdArgs = parts[1]
	}
	return NewCommand(cmdName, cmdArgs), nil
}
