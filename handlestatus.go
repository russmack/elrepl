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
	h.HandlerFunc = func(cmd *Command) string {
		argParts := strings.Split(cmd.Args, " ")
		//arg := "_status?pretty"
		index := ""
		//urlString := ""
		res := ""
		if len(argParts) == 0 {
			u := new(url.URL)
			u.Scheme = "http"
			u.Host = server.host + ":" + server.port
			u.Path = "_status"
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
			u.Path = index + "/" + "_status"
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
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
