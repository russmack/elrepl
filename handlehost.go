package main

import ()

func init() {
	h := NewHandler()
	h.CommandName = "host"
	h.CommandPattern = "(host)(( )(.*))"
	h.Usage = "host [hostAddress]"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		p := ParseMap{}

		switch cmd.Args[0] {
		case "/?":
			return p, false
		case "":
			if len(cmd.Args) == 1 { // get host
				return p, true
			} else {
				return p, false
			}
		default:
			if len(cmd.Args) == 1 { // set host
				p["host"] = cmd.Args[0]
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
		hostAddress, ok := p["host"]
		if !ok {
			return "Host: " + server.host
		} else {
			server.host = hostAddress
			return "Set host: " + hostAddress
		}
	}
	HandlerRegistry[h.CommandName] = h
}
