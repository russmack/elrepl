package main

import (
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "port"
	h.CommandPattern = "(port)( )(.*)"
	h.Usage = "port [portNumber]"
	h.CommandParser = func(cmd *Command) (map[string]string, bool) {
		if cmd.Args == "/?" {
			return nil, false
		}
		argParts := strings.Split(cmd.Args, " ")
		if len(argParts) > 1 {
			return nil, false
		}
		if len(argParts) == 1 && argParts[0] == "" {
			return nil, true
		} else {
			m := make(map[string]string)
			m["port"] = argParts[0]
			return m, true
		}
	}
	h.HandlerFunc = func(cmd *Command) string {
		m, ok := h.CommandParser(cmd)
		if !ok {
			return usageMessage(h.Usage)
		}
		port, ok := m["port"]
		if !ok {
			return "Host: " + server.port
		} else {
			server.port = port
			return "Set port: " + port
		}
	}
	HandlerRegistry[h.CommandName] = h
}
