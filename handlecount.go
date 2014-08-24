package main

import (
	"fmt"
	"net/url"
)

func init() {
	h := NewHandler()
	h.CommandName = "count"
	h.CommandPattern = "(count)( )(.*)"
	h.Usage = "count [/|indexName] [type]"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		p := ParseMap{}
		p["scheme"] = "http"
		p["host"] = server.host
		p["port"] = server.port
		p["endpoint"] = "_count"

		switch cmd.Args[0] {
		case "/?":
			return p, false
		case "":
			if len(cmd.Args) == 1 {
				return p, true
			} else {
				return p, false
			}
		default:
			if len(cmd.Args) == 1 {
				p["index"] = cmd.Args[0]
				return p, true
			} else if len(cmd.Args) == 2 {
				p["index"] = cmd.Args[0]
				p["type"] = cmd.Args[1]
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
		docType, ok := p["type"]
		if ok {
			index += docType + "/"
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