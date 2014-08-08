package main

import (
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "index"
	h.CommandPattern = "(index)(( )(.*))"
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
			m["index"] = argParts[0]
			return m, true
		}
	}
	h.HandlerFunc = func(cmd *Command) string {
		m, ok := h.CommandParser(cmd)
		if !ok {
			return "Usage: index [indexName]"
		}
		indexName, ok := m["index"]
		if !ok {
			return "Index: " + server.index
		} else {
			server.index = indexName
			return "Set index: " + indexName
		}
	}
	HandlerRegistry[h.CommandName] = h
}
