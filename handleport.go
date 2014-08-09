package main

import (
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "port"
	h.CommandPattern = "(port)( )(.*)"
	h.Usage = "port [portNumber]"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		argParts := strings.Split(cmd.Args, " ")
		p := ParseMap{}

		switch argParts[0] {
		case "/?":
			return p, false
		case "":
			if len(argParts) == 1 { // get port
				return p, true
			} else {
				return p, false
			}
		default:
			if len(argParts) == 1 { // set port
				p["port"] = argParts[0]
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
		port, ok := p["port"]
		if !ok {
			return "Host: " + server.port
		} else {
			server.port = port
			return "Set port: " + port
		}
	}
	HandlerRegistry[h.CommandName] = h
}
