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
	h.CommandParser = func(cmd *Command) map[string]string {
		argParts := strings.Split(cmd.Args, " ")
		m := make(map[string]string)
		m["scheme"] = "http"
		m["host"] = server.host
		m["port"] = server.port
		m["endpoint"] = "_status"
		if argParts[0] != "/" {
			m["index"] = argParts[0]
		}
		return m
	}
	h.HandlerFunc = func(cmd *Command) string {
		m := h.CommandParser(cmd)
		u := new(url.URL)
		u.Scheme = m["scheme"]
		u.Host = m["host"] + ":" + m["port"]
		index, ok := m["index"]
		if ok {
			index += "/"
		}
		u.Path = index + m["endpoint"]
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
