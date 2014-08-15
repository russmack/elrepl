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
	//parts := strings.SplitN(entry, " ", 2)
	parts := strings.Split(entry, " ")
	cmdName := parts[0]
	cmdArgs := parts[1:]
	if len(cmdArgs) == 0 {
		cmdArgs = append(cmdArgs, "")
	}
	return NewCommand(cmdName, cmdArgs), nil
}
