package main

import (
	"fmt"
	"net/url"
)

func init() {
	h := NewHandler()
	h.CommandName = "version"
	h.CommandPattern = "(version)(( )(.*))"
	h.CommandParser = func(cmd *Command) map[string]string {
		m := make(map[string]string)
		m["scheme"] = "http"
		m["host"] = server.host
		m["port"] = server.port
		return m
	}
	h.HandlerFunc = func(cmd *Command) string {
		m := h.CommandParser(cmd)
		u := new(url.URL)
		u.Scheme = m["scheme"]
		u.Host = m["host"] + ":" + m["port"]
		fmt.Println("Request:", u)
		res, err := getHttpResource(u.String())
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
