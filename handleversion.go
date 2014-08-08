package main

import (
	"fmt"
	"net/url"
)

func init() {
	h := NewHandler()
	h.CommandName = "version"
	h.CommandPattern = "(version)(( )(.*))"
	h.Usage = "version"
	h.CommandParser = func(cmd *Command) (map[string]string, bool) {
		if cmd.Args == "/?" {
			return nil, false
		}
		if cmd.Args != "" {
			return nil, false
		}
		m := make(map[string]string)
		m["scheme"] = "http"
		m["host"] = server.host
		m["port"] = server.port
		return m, true
	}
	h.HandlerFunc = func(cmd *Command) string {
		m, ok := h.CommandParser(cmd)
		if !ok {
			return usageMessage(h.Usage)
		}
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
