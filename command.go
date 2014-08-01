package main

import (
	"regexp"
)

type Command struct {
	Name    string
	Args    string
	Pattern *regexp.Regexp
}

//func NewCommand(name string, args string, pattern *regexp.Regexp) *Command {
func NewCommand(name string, args string) *Command {
	//return &Command{Name: name, Args: args, Pattern: pattern}
	return &Command{Name: name, Args: args}
}

/*
func NewCommand(name string, regexp *regexp) *Command {
	return &Command{Name: name, Pattern: regexp}
}
*/
