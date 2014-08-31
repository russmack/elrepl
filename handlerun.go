package main

import (
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "run"
	h.CommandPattern = "(run)( )(.*)"
	h.Usage = "run"
	h.HandlerFunc = func(cmd *Command) string {
		verb := strings.ToLower(cmd.Tokens[0])
		cmdParser := NewCommandParser()
		newCmd, err := cmdParser.Parse(loadedRequest.request)
		if err != nil {
			return "Unable to parse loaded query for run command."
		}
		dispatcher := NewDispatcher()
		switch verb {
		case "post":
			resp := dispatcher.Dispatch(newCmd)
			return resp
		case "put":
			resp := dispatcher.Dispatch(newCmd)
			return resp
		case "get":
			resp := dispatcher.Dispatch(newCmd)
			return resp
		}
		return "Unable to run loaded query."
	}
	HandlerRegistry[h.CommandName] = h
}
