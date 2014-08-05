package main

import (
	"fmt"
	//"net/url"
	"strings"
)

func init() {
	h := NewHandler()
	h.CommandName = "doc"
	h.CommandPattern = "(doc)(( )(.*))"
	h.HandlerFunc = func(cmd *Command) string {
		argParts := strings.SplitN(cmd.Args, " ", 3)

		arg := "_mget?pretty"
		url := ""
		res := ""
		fmt.Println("ARGS:", argParts)
		if argParts[0] == "get" {
			indexName := argParts[1]
			ids := argParts[2]
			//url = fmt.Sprintf("http://%s:%s", server.host, server.port)
			//indexName := argParts[1]
			//aliasName := argParts[2]
			// "ids" : ["1", "1"]
			url = "post " + indexName + "/" + arg + " " + "{\"ids\": " + ids + " }"
			cmdParser := NewCommandParser()
			newCmd, err := cmdParser.Parse(url)
			if err != nil {
				return err.Error()
			}
			dispatcher := NewDispatcher()
			res = dispatcher.Dispatch(newCmd)
		} else if argParts[0] == "delete" {
			//curl -XDELETE 'http://localhost:9200/twitter/tweet/1'
			indexName := argParts[1]
			delArgs := strings.Split(argParts[2], " ")
			typeName := delArgs[0]
			id := delArgs[1]
			url = fmt.Sprintf("http://%s:%s/%s/%s/%s", server.host, server.port, indexName, typeName, id)
			fmt.Println("Request:", url)
			delRes, err := deleteHttpResource(url)
			if err != nil {
				return err.Error()
			}
			res = delRes
		}
		//fmt.Println("Request:", url)
		//res, err := getHttpResource(url)
		//if err != nil {
		//	return err.Error()
		//}
		return res
	}
	HandlerRegistry[h.CommandName] = h
}
