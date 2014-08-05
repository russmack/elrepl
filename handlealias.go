package main

import (
	"fmt"
	"net/url"
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "alias"
	h.CommandPattern = "(alias)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		argParts := strings.Split(cmd.Args, " ")
		arg := "_aliases?pretty"
		index := ""
		urlString := ""
		res := ""
		if len(argParts) == 0 {
			u := new(url.URL)
			u.Scheme = "http"
			u.Host = server.host + ":" + server.port
			u.Path = "_aliases"
			q := u.Query()
			q.Add("pretty", "true")
			u.RawQuery = q.Encode()
			fmt.Println("Request:", u)
			getRes, err := getHttpResource(u.String())
			if err != nil {
				return err.Error()
			}
			res = getRes
		} else if len(argParts) == 1 {
			if argParts[0] == "index" { // Override session var.
				index = server.index
			} else {
				index = argParts[0]
			}

			u := new(url.URL)
			u.Scheme = "http"
			u.Host = server.host + ":" + server.port
			u.Path = index + "/" + "_aliases"
			q := u.Query()
			q.Add("pretty", "true")
			u.RawQuery = q.Encode()

			//urlString = fmt.Sprintf("http://%s:%s/%s/%s", server.host, server.port, index, arg)
			fmt.Println("Request:", u)
			getRes, err := getHttpResource(u.String())
			if err != nil {
				return err.Error()
			}
			res = getRes
		} else {
			switch argParts[0] {
			case "create":
				//curl -XPOST "http://10.1.1.12:9200/_aliases" -d '{ "actions": [ { "add": { "index": "podcasts-2014-07-29-001", "alias": "podcastsupdater" } } ] }'
				//post _search?pretty { "query": { "term": { "director": "scott" } } }
				indexName := argParts[1]
				aliasName := argParts[2]
				urlString = "post " + arg + " " + "{\"actions\": [ { \"add\": { \"index\": \"" + indexName + "\", \"alias\": \"" + aliasName + "\" } } ] }"
				cmdParser := NewCommandParser()
				newCmd, err := cmdParser.Parse(urlString)
				if err != nil {
					return err.Error()
				}
				dispatcher := NewDispatcher()
				res = dispatcher.Dispatch(newCmd)
			case "remove":
				//curl -XPOST "http://10.1.1.12:9200/_aliases" -d '{ "actions": [ { "remove": { "index": "podcasts-2014-05-07-0103", "alias": "podcastsupdater" } } ] }'
				indexName := argParts[1]
				aliasName := argParts[2]
				urlString = "post " + arg + " " + "{\"actions\": [ { \"remove\": { \"index\": \"" + indexName + "\", \"alias\": \"" + aliasName + "\" } } ] }"
				cmdParser := NewCommandParser()
				newCmd, err := cmdParser.Parse(urlString)
				if err != nil {
					return err.Error()
				}
				dispatcher := NewDispatcher()
				res = dispatcher.Dispatch(newCmd)
			default:
				res = "Commmand args not recognised."
			}
		}

		return res
	}
	HandlerRegistry[h.CommandName] = h
}
