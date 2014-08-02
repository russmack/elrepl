package main

func init() {
	h := NewHandler()
	h.CommandName = "index"
	h.CommandPattern = "(index)(( )(.*))"
	h.HandlerFunc = func(cmd *Command) string {
		if cmd.Args == "" {
			return "Index: " + server.index
		} else {
			arg := cmd.Args
			server.index = arg
			return "Set index: " + arg
		}
	}
	HandlerRegistry[h.CommandName] = h
}
