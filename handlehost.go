package main

import (
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "host"
	h.CommandPattern = "(host)(( )(.*))"
	h.Usage = "host [hostAddress]"
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
			m["host"] = argParts[0]
			return m, true
		}
	}
	h.HandlerFunc = func(cmd *Command) string {
		m, ok := h.CommandParser(cmd)
		if !ok {
			return usageMessage(h.Usage)
		}
		hostAddress, ok := m["host"]
		if !ok {
			return "Host: " + server.host
		} else {
			server.host = hostAddress
			return "Set host: " + hostAddress
		}
	}
	HandlerRegistry[h.CommandName] = h
}
