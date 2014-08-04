package main

import (
	"fmt"
)

func init() {
	h := NewHandler()
	h.CommandName = "alias"
	h.CommandPattern = "(alias)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		//arg := cmd.Args
		arg := "_aliases?pretty"

		url := ""
		if server.index == "" {
			url = fmt.Sprintf("http://%s:%s/%s", server.host, server.port, arg)
		} else {
			url = fmt.Sprintf("http://%s:%s/%s/%s", server.host, server.port, server.index, arg)
		}

		fmt.Println("Request:", url)
		res, err := getHttpResource(url)
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
