package main

import (
	"fmt"
	"net/url"
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "status"
	h.CommandPattern = "(status)( )(.*)"
	h.Usage = "status [/|indexName]"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		argParts := strings.Split(cmd.Args, " ")
		p := ParseMap{}
		p["scheme"] = "http"
		p["host"] = server.host
		p["port"] = server.port
		p["endpoint"] = "_status"

		switch argParts[0] {
		case "/?":
			return p, false
		case "":
			if len(argParts) == 1 {
				return p, true
			} else {
				return p, false
			}
		default:
			if len(argParts) == 1 {
				p["index"] = argParts[0]
				return p, true
			} else {
				return p, false
			}
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
		index, ok := p["index"]
		if ok {
			index += "/"
		}
		u.Path = index + p["endpoint"]
		q := u.Query()
		q.Add("pretty", "true")
		u.RawQuery = q.Encode()
		fmt.Println("Request:", u)
		res, err := getHttpResource(u.String())
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
