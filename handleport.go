package main

import (
	"regexp"
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "port"
	h.CommandPattern = "(port)( )(.*)"
	h.Usage = "port [portNumber]"
	h.CommandParser = func(cmd *Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Get port
			regexp.MustCompile(`^port$`): func(s []string) (string, bool) {
				r := "Port: " + server.port
				return r, true
			},
			// Port help
			regexp.MustCompile(`^port /\?$`): func(s []string) (string, bool) {
				return "", false
			},
			// Set port
			regexp.MustCompile(`^port ([a-zA-Z0-9\.]+)$`): func(s []string) (string, bool) {
				server.port = s[1]
				r := "Set port: " + server.port
				return r, true
			},
		}
		r, ok := h.Tokenize(strings.TrimSpace(cmd.Instruction), pattFn)
		return r, ok
	}
	h.HandlerFunc = func(cmd *Command) string {
		r, ok := h.CommandParser(cmd)
		if !ok {
			return r + "\n\n" + usageMessage(h.Usage)
		}
		return r
	}
	HandlerRegistry[h.CommandName] = h
}
