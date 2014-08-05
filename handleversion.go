package main

import (
	"fmt"
	"net/url"
)

func init() {
	h := NewHandler()
	h.CommandName = "version"
	h.CommandPattern = "(version)(( )(.*))"
	h.HandlerFunc = func(cmd *Command) string {
		u := new(url.URL)
		u.Scheme = "http"
		u.Host = server.host + ":" + server.port
		fmt.Println("Request:", u)
		res, err := getHttpResource(u.String())
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
