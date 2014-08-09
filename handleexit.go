package main

import (
	"fmt"
	"os"
)

func init() {
	h := NewHandler()
	h.CommandName = "exit"
	h.CommandPattern = "(exit)( )(.*)"
	h.Usage = "exit"
	h.HandlerFunc = func(cmd *Command) string {
		fmt.Println("Bye.")
		os.Exit(0)
		return ""
	}
	HandlerRegistry[h.CommandName] = h
}
