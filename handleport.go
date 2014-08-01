package main

func init() {
	h := NewHandler()
	h.CommandName = "port"
	h.CommandPattern = "(port)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		if cmd.Args == "" {
			return "Server port: " + server.port
		} else {
			arg := cmd.Args
			server.port = arg
			return "Set server port: " + arg
		}
	}
	HandlerRegistry[h.CommandName] = h
}
