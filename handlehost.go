package main

import (
	"regexp"
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "host"
	h.CommandPattern = "(host)(( )(.*))"
	h.Usage = "host [hostAddress]"
	h.CommandParser = func(cmd *Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Get host
			regexp.MustCompile(`^host$`): func(s []string) (string, bool) {
				r := "Host: " + server.host
				return r, true
			},
			// Host help
			regexp.MustCompile(`^host /\?$`): func(s []string) (string, bool) {
				return "", false
			},
			// Set help
			regexp.MustCompile(`^host ([a-zA-Z0-9\.]+)$`): func(s []string) (string, bool) {
				server.host = s[1]
				r := "Set host: " + server.host
				return r, true
			},
		}
		r, ok := h.Tokenize(strings.TrimSpace(cmd.Instruction), pattFn)
		return r, ok
	}
	h.HandlerFunc = func(cmd *Command) string {
		r, ok := h.CommandParser(cmd)
		if !ok {
			return usageMessage(h.Usage)
		}
		return r
	}
	HandlerRegistry[h.CommandName] = h
}
