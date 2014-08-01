package main

//import (
//"fmt"
//"github.com/russmack/elrepl/types"
//)

func init() {
	h := NewHandler()
	h.CommandName = "version"
	h.CommandPattern = "(version)(( )(.*))"
	h.HandlerFunc = func(cmd *Command) string {
		return `
	elRepl version 0.1
	`
	}
	HandlerRegistry[h.CommandName] = h
}
