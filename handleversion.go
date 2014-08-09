package main

import (
	"fmt"
	"net/url"
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "version"
	h.CommandPattern = "(version)(( )(.*))"
	h.Usage = "version"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		argParts := strings.Split(cmd.Args, " ")
		p := ParseMap{}
		p["scheme"] = "http"
		p["host"] = server.host
		p["port"] = server.port

		switch argParts[0] {
		case "/?":
			return p, false
		case "":
			return p, true
		default:
			return p, false
		}
	}
	h.HandlerFunc = func(cmd *Command) string {
		p, ok := h.CommandParser(cmd)
		if !ok {
			return usageMessage(h.Usage)
		}
		u := new(url.URL)
		u.Scheme = p["scheme"]
		u.Host = p["host"] + ":" + p["port"]
		fmt.Println("Request:", u)
		res, err := getHttpResource(u.String())
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
