package main

import ()

func init() {
	h := NewHandler()
	h.CommandName = "host"
	h.CommandPattern = "(host)(( )(.*))"
	h.HandlerFunc = func(cmd *Command) string {
		if cmd.Args == "" {
			return "Server host: " + server.host
		} else {
			arg := cmd.Args
			server.host = arg
			return "Set server host: " + arg
		}
	}
	HandlerRegistry[h.CommandName] = h
}
