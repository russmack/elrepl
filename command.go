package main

import (
	"regexp"
)

type Command struct {
	Name    string
	Args    string
	Pattern *regexp.Regexp
}

func NewCommand(name string, args string) *Command {
	return &Command{Name: name, Args: args}
}
