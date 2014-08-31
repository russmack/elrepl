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
	bodyStart := strings.Index(entry, "{")
	//parts := strings.SplitN(entry, "{", 2)
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
	//if len(parts) > 0 {
	//instr = parts[0]
	//if len(parts) > 1 {
	//	body = parts[1]
	//}
	cmdTokens := strings.Split(instr, " ")
	cmdName := cmdTokens[0]
	//cmdArgs := tokens[1:]
	//if len(cmdArgs) == 0 {
	//	cmdArgs = append(cmdArgs, "")
	//}
	return NewCommand(instr, body, cmdName, cmdTokens), nil
}
