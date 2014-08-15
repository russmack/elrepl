package main

import (
	"errors"
	"fmt"
	"net/url"
)

func init() {
	h := NewHandler()
	h.CommandName = "alias"
	h.CommandPattern = "(alias)( )(.*)"
	h.Usage = "alias (create host port indexName aliasName) | (remove host port indexName aliasName) | (move host port fromIndex toIndex aliasName)"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		p := ParseMap{}
		p["scheme"] = "http"
		p["endpoint"] = "_aliases"
		switch cmd.Args[0] {
		case "/?":
			return p, false
		case "":
			if len(cmd.Args) == 1 { // list all
				// Session vars for read ops only.
				p["action"] = "list"
				p["host"] = server.host
				p["port"] = server.port
			} else {
				return p, false
			}
		case "create":
			p["action"] = "create"
			p["host"] = cmd.Args[1]
			p["port"] = cmd.Args[2]
			p["index"] = cmd.Args[3]
			p["alias"] = cmd.Args[4]
		case "remove":
			p["action"] = "remove"
			p["host"] = cmd.Args[1]
			p["port"] = cmd.Args[2]
			p["index"] = cmd.Args[3]
			p["alias"] = cmd.Args[4]
		case "move":
			p["action"] = "move"
			p["host"] = cmd.Args[1]
			p["port"] = cmd.Args[2]
			p["fromIndex"] = cmd.Args[3]
			p["toIndex"] = cmd.Args[4]
			p["alias"] = cmd.Args[5]
		default:
			if len(cmd.Args) == 1 { // list index
				p["action"] = "list"
				p["host"] = server.host
				p["port"] = server.port
				p["index"] = cmd.Args[0]
			} else {
				return p, false
			}
		}
		return p, true
	}
	h.HandlerFunc = func(cmd *Command) string {
		p, ok := h.CommandParser(cmd)
		if !ok {
			return usageMessage(h.Usage)
		}

		res := ""
		switch p["action"] {
		case "list":
			u := new(url.URL)
			u.Scheme = p["scheme"]
			u.Host = p["host"] + ":" + p["port"]
			if index, ok := p["index"]; ok {
				u.Path = index + "/" + p["endpoint"]
			} else {
				u.Path = p["endpoint"]
			}
			q := u.Query()
			q.Add("pretty", "true")
			u.RawQuery = q.Encode()
			fmt.Println("Request:", u)
			err := errors.New("")
			res, err = getHttpResource(u.String())
			if err != nil {
				return err.Error()
			}
		case "create":
			//curl -XPOST "http://10.1.1.12:9200/_aliases" -d '{ "actions": [ { "add": { "index": "podcasts-2014-07-29-001", "alias": "podcastsupdater" } } ] }'
			//post _alias?pretty { "actions": [ { "add": { "index": "podcasts-2014-07-29-001", "alias": "podcastsupdater" } } ] }
			indexName := p["index"]
			aliasName := p["alias"]
			urlString := "post " + p["endpoint"] + " " + "{\"actions\": [ { \"add\": { \"index\": \"" + indexName + "\", \"alias\": \"" + aliasName + "\" } } ] }"
			cmdParser := NewCommandParser()
			newCmd, err := cmdParser.Parse(urlString)
			if err != nil {
				return err.Error()
			}
			dispatcher := NewDispatcher()
			res = dispatcher.Dispatch(newCmd)
		case "remove":
			//curl -XPOST "http://10.1.1.12:9200/_aliases" -d '{ "actions": [ { "remove": { "index": "podcasts-2014-05-07-0103", "alias": "podcastsupdater" } } ] }'
			indexName := p["index"]
			aliasName := p["alias"]
			urlString := "post " + p["endpoint"] + " " + "{\"actions\": [ { \"remove\": { \"index\": \"" + indexName + "\", \"alias\": \"" + aliasName + "\" } } ] }"
			cmdParser := NewCommandParser()
			newCmd, err := cmdParser.Parse(urlString)
			if err != nil {
				return err.Error()
			}
			dispatcher := NewDispatcher()
			res = dispatcher.Dispatch(newCmd)
		case "move":
			fromIndexName := p["fromIndex"]
			toIndexName := p["toIndex"]
			aliasName := p["alias"]
			urlString := "post " + p["endpoint"] + " " + "{ \"actions\": [ { \"remove\": { \"alias\": \"" + aliasName + "\", \"index\": \"" + fromIndexName + "\" }}, { \"add\": { \"alias\": \"" + aliasName + "\", \"index\": \"" + toIndexName + "\" } } ] }"
			cmdParser := NewCommandParser()
			newCmd, err := cmdParser.Parse(urlString)
			if err != nil {
				return err.Error()
			}
			dispatcher := NewDispatcher()
			res = dispatcher.Dispatch(newCmd)
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
