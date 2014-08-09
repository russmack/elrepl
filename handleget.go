package main

import (
	"fmt"
	"net/url"
)

func init() {
	h := NewHandler()
	h.CommandName = "get"
	h.CommandPattern = "(get)( )(.*)"
	h.Usage = "get url"
	h.HandlerFunc = func(cmd *Command) string {
		arg := cmd.Args

		u := new(url.URL)
		newUrl, err := u.Parse(arg)
		if err != nil {
			return "Unable to parse url: " + err.Error()
		}

		fmt.Println("Request:", newUrl)
		res, err := getHttpResource(newUrl.String())
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
