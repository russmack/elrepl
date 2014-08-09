package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "alias"
	h.CommandPattern = "(alias)( )(.*)"
	h.Usage = "alias (create indexName aliasName) | (remove indexName aliasName) | (move fromIndex toIndex aliasName)"
	h.CommandParse = func(cmd *Command) (ParseMap, bool) {
		argParts := strings.Split(cmd.Args, " ")
		p := ParseMap{}
		p["scheme"] = "http"
		p["host"] = server.host
		p["port"] = server.port
		p["endpoint"] = "_aliases"

		switch argParts[0] {
		case "/?":
			return p, false
		case "":
			if len(argParts) == 1 { // list all
				p["action"] = "list"
				return p, true
			} else {
				return p, false
			}
		case "create":
			p["action"] = "create"
			p["index"] = argParts[1]
			p["alias"] = argParts[2]
		case "remove":
			p["action"] = "remove"
			p["index"] = argParts[1]
			p["alias"] = argParts[2]
		case "move":
			p["action"] = "move"
			p["fromIndex"] = argParts[1]
			p["toIndex"] = argParts[2]
			p["alias"] = argParts[3]
		default:
			if len(argParts) == 1 { // list index
				p["action"] = "list"
				p["index"] = argParts[0]
			} else {
				return p, false
			}
		}
		return p, true
	}
	h.HandlerFunc = func(cmd *Command) string {
		p, ok := h.CommandParse(cmd)
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
