package main

import ()

func init() {
	h := NewHandler()
	h.CommandName = "index"
	h.CommandPattern = "(index)(( )(.*))"
	h.Usage = "index [indexName]"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		p := ParseMap{}

		switch cmd.Args[0] {
		case "/?":
			return p, false
		case "":
			if len(cmd.Args) == 1 { // get index
				return p, true
			} else {
				return p, false
			}
		default:
			if len(cmd.Args) == 1 { // set index
				p["index"] = cmd.Args[0]
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
