package main

import (
	"errors"
	"fmt"
	"net/url"
)

func init() {
	h := NewHandler()
	h.CommandName = "index"
	h.CommandPattern = "(index)(( )(.*))"
	h.Usage = "index [ indexName | (create host port indexName) ]"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		p := ParseMap{}
		p["scheme"] = "http"

		switch cmd.Args[0] {
		case "/?":
			return p, false
		case "":
			if len(cmd.Args) == 1 { // get index
				return p, true
			} else {
				return p, false
			}
		case "create":
			if len(cmd.Args) == 4 { // set index
				p["action"] = cmd.Args[0]
				p["host"] = cmd.Args[1]
				p["port"] = cmd.Args[2]
				p["index"] = cmd.Args[3]
				return p, true
			} else {
				return p, false
			}
		default:
			if len(cmd.Args) == 1 { // set index
				p["index"] = cmd.Args[0]
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
		action, actionOk := p["action"]
		indexName, indexOk := p["index"]
		if !indexOk {
			return "Index: " + server.index
		}
		if !actionOk {
			server.index = indexName
			return "Set index: " + indexName
		}
		if action == "create" {
			// Create new index
			//curl -XPUT "http://10.1.1.12:9200/podcasts-2014-07-23-1930/"
			u := new(url.URL)
			u.Scheme = p["scheme"]
			u.Host = p["host"] + ":" + p["port"]
			u.Path = indexName
			q := u.Query()
			q.Add("pretty", "true")
			u.RawQuery = q.Encode()
			fmt.Println("Request:", u)
			err := errors.New("")
			res, err := putHttpResource(u.String(), "")
			if err != nil {
				return err.Error()
			}
			return res
		}
		return usageMessage(h.Usage)
	}
	HandlerRegistry[h.CommandName] = h
}
