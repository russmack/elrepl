package main

import (
	"fmt"
)

func init() {
	h := NewHandler()
	h.CommandName = "version"
	h.CommandPattern = "(version)(( )(.*))"
	h.HandlerFunc = func(cmd *Command) string {
		url := fmt.Sprintf("http://%s:%s", server.host, server.port)
		fmt.Println("Request:", url)
		res, err := getHttpResource(url)
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
