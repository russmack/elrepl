package main

import (
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "run"
	h.CommandPattern = "(run)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		loadedParts := strings.SplitN(loadedRequest.request, "\n", 2)

		loadedCmdParts := strings.SplitN(loadedParts[0], " ", 2)
		loadedCmd := loadedCmdParts[0]
		//loadedArgs := loadedCmdParts[1]

		//loadedQuery := loadedParts[1]

		cmdParser := NewCommandParser()
		newCmd, err := cmdParser.Parse(loadedRequest.request)
		if err != nil {
			return "Unable to parse loaded query for run command."
		}
		dispatcher := NewDispatcher()
		if strings.ToLower(loadedCmd) == "post" {
			//resp := HandlerRegistry[newCmd.Name]
			resp := dispatcher.Dispatch(newCmd)
			return resp
		} else if strings.ToLower(loadedCmd) == "put" {
			//resp := handlePut(newCmd)
			resp := dispatcher.Dispatch(newCmd)
			return resp
		} else if strings.ToLower(loadedCmd) == "get" {
			//resp := handleGet(newCmd)
			resp := dispatcher.Dispatch(newCmd)
			return resp
		}
		return "Unable to run loaded query."
	}
	HandlerRegistry[h.CommandName] = h
}
