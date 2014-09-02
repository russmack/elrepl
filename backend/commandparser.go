package backend

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
	bodyStart := strings.Index(entry, "{")
	instr := ""
	body := ""
	if bodyStart == -1 {
		instr = entry
	} else {
		instr = entry[:bodyStart]
	}
	if bodyStart > -1 {
		body = entry[bodyStart:]
	}
	cmdTokens := strings.Split(instr, " ")
	cmdName := cmdTokens[0]
	return NewCommand(instr, body, cmdName, cmdTokens), nil
}
