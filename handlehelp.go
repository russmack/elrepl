package main

import ()

func init() {
	h := NewHandler()
	h.CommandName = "help"
	h.CommandPattern = "(help)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		cmdList := ""

		for k, _ := range HandlerRegistry {
			cmdList += k + "\n"
		}

		return `
Help
----
Commands:
` + cmdList + `
eg:
host localhost
port 9200
index movies
get _search?q=title:thx1138
`
	}
	HandlerRegistry[h.CommandName] = h
}
