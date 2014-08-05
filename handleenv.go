package main

func init() {
	h := NewHandler()
	h.CommandName = "env"
	h.CommandPattern = "(env)(( )(.*))"
	h.HandlerFunc = func(cmd *Command) string {
		return `
	elRepl version 0.1
	`
	}
	HandlerRegistry[h.CommandName] = h
}
