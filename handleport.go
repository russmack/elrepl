package main

import ()

func init() {
	h := NewHandler()
	h.CommandName = "port"
	h.CommandPattern = "(port)( )(.*)"
	h.Usage = "port [portNumber]"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		p := ParseMap{}

		switch cmd.Args[0] {
		case "/?":
			return p, false
		case "":
			if len(cmd.Args) == 1 { // get port
				return p, true
			} else {
				return p, false
			}
		default:
			if len(cmd.Args) == 1 { // set port
				p["port"] = cmd.Args[0]
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
