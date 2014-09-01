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
		//queryHost := cmd.Tokens[1]
		//queryPort := cmd.Tokens[2]
		//queryArgs := cmd.Tokens[3]

		arg := cmd.Tokens[1]
		u := new(url.URL)
		newUrl, err := u.Parse(arg)
		if err != nil {
			return "Unable to parse url: " + err.Error()
		}
		fmt.Println("Request: get", newUrl)
		res, err := getHttpResource(newUrl.String())
		if err != nil {
			return err.Error()
		}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
