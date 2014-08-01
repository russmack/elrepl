package main

import (
	"regexp"
	"strings"
)

type CommandParser struct {
	commandRegexpMap map[string]*regexp.Regexp
}

func NewCommandParser() *CommandParser {
	p := CommandParser{}
	//cmds, err := p.initCommands()
	//if err != nil {
	//	panic(err)
	//}
	//p.commandRegexpMap = cmds
	return &p
}

func (p *CommandParser) Parse(entry string) (*Command, error) {
	parts := strings.SplitN(entry, " ", 2)
	cmdName := parts[0]
	cmdArgs := ""
	if len(parts) > 1 {
		cmdArgs = parts[1]
	}
	//cmdRe, ok := p.commandRegexpMap[cmdName]
	//if ok {
	//	return NewCommand(cmdName, cmdArgs, cmdRe), nil
	//} else {
	//	return nil, nil
	//}
	return NewCommand(cmdName, cmdArgs), nil
}

/*
func (p *CommandParser) initCommands() (map[string]*regexp.Regexp, error) {
	commandPatternMap := make(map[string]string)
	commandRegexpMap := make(map[string]*regexp.Regexp)

	commandPatternMap[Commands.Post] = "((?i)post(?-i))( )(.*)"
	commandPatternMap[Commands.Reindex] = "(reindex)( )(.*)"

	commandPatternMap[Commands.Put] = "(put)( )(.*)"
	commandPatternMap[Commands.Get] = "(get)( )(.*)"
	commandPatternMap[Commands.Run] = "(run)( )(.*)"
	commandPatternMap[Commands.Load] = "(load)( )(.*)"
	commandPatternMap[Commands.Log] = "(log)( )(.*)"
	commandPatternMap[Commands.Port] = "(port)( )(.*)"
	commandPatternMap[Commands.Host] = "(host)(( )(.*))"
	commandPatternMap[Commands.Dir] = "(dir)( )(.*)"
	commandPatternMap[Commands.Exit] = "(exit)( )(.*)"
	commandPatternMap[Commands.Help] = "(help)( )(.*)"
	commandPatternMap[Commands.Version] = "(version)(( )(.*))"

	for k, v := range commandPatternMap {
		r, err := regexp.Compile(v)
		if err != nil {
			return nil, err
		}
		commandRegexpMap[k] = r
	}
	return commandRegexpMap, nil
}
*/
