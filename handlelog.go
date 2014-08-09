package main

import (
	"strconv"
)

func init() {
	h := NewHandler()
	h.CommandName = "log"
	h.CommandPattern = "(log)( )(.*)"
	h.Usage = "log"
	h.HandlerFunc = func(cmd *Command) string {
		logLevel = 1
		return "Logging level set to: " + strconv.Itoa(logLevel)
	}
	HandlerRegistry[h.CommandName] = h
}
