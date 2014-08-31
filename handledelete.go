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
	// TODO: usage
	h.HandlerFunc = func(cmd *Command) string {
		arg := cmd.Args[0]
		u := new(url.URL)
		newUrl, err := u.Parse(arg)
		if err != nil {
			return "Unable to parse url: " + err.Error()
		}
		fmt.Println("Request:", newUrl)
		res, err := deleteHttpResource(newUrl.String())
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
