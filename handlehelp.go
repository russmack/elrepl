package main

import (
	"sort"
)

func init() {
	h := NewHandler()
	h.CommandName = "help"
	h.CommandPattern = "(help)( )(.*)"
	h.Usage = "help"
	h.HandlerFunc = func(cmd *Command) string {
		cmdList := ""
		keys := []string{}
		for k, _ := range HandlerRegistry {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, s := range keys {
			cmdList += s + ": " + HandlerRegistry[s].Usage + "\n"
		}
		return `
	Help
	----


Write commands require host and port be specified to avoid slip-ups.
Read commands can omit host and port, and session settings will be used.


Commands
--------
` + cmdList + `
`
	}
	HandlerRegistry[h.CommandName] = h
}
