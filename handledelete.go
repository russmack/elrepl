package main

import (
	"fmt"
	"net/url"
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "delete"
	h.CommandPattern = "(delete)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		argParts := strings.Split(cmd.Args, " ")

		// TODO: This handler is not implemented.

		u := new(url.URL)
		u.Scheme = "http"
		u.Host = server.host + ":" + server.port
		u.Path = server.index + "/" + "_aliases"
		u.RawQuery = u.Query().Add("pretty", "true")
		//if server.index == "" {
		//	url = fmt.Sprintf("http://%s:%s/%s", server.host, server.port, arg)
		//} else {
		//url = fmt.Sprintf("http://%s:%s/%s/%s", server.host, server.port, server.index, arg)
		//}

		fmt.Println("Request:", u)
		res, err := deleteHttpResource(u.String())
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
