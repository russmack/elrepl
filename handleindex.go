package main

import (
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "index"
	h.CommandPattern = "(index)(( )(.*))"
	h.Usage = "index [indexName]"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		argParts := strings.Split(cmd.Args, " ")
		p := ParseMap{}

		switch argParts[0] {
		case "/?":
			return p, false
		case "":
			if len(argParts) == 1 { // get index
				return p, true
			} else {
				return p, false
			}
		default:
			if len(argParts) == 1 { // set index
				p["index"] = argParts[0]
				return p, true
			} else {
				return p, false
			}
		}
	}
	h.HandlerFunc = func(cmd *Command) string {
		p, ok := h.CommandParser(cmd)
		if !ok {
			return usageMessage(h.Usage)
		}
		indexName, ok := p["index"]
		if !ok {
			return "Index: " + server.index
		} else {
			server.index = indexName
			return "Set index: " + indexName
		}
	}
	HandlerRegistry[h.CommandName] = h
}
