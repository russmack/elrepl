package main

import ()

func init() {
	h := NewHandler()
	h.CommandName = "help"
	h.CommandPattern = "(help)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		return `
	Help
	----
	Commands:
	  eg:
	  host localhost
	  port 9200
	  index movies
	  get _search?q=title:thx1138
	`
	}
	HandlerRegistry[h.CommandName] = h
}
