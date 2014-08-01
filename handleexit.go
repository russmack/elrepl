package main

import (
	"fmt"
	//"github.com/russmack/elrepl/types"
	"os"
)

func init() {
	h := NewHandler()
	h.CommandName = "exit"
	h.CommandPattern = "(exit)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		fmt.Println("Bye.")
		os.Exit(0)
		return ""
	}
	HandlerRegistry[h.CommandName] = h
}
