package main

import (
	"io/ioutil"
)

func init() {
	h := NewHandler()
	h.CommandName = "load"
	h.CommandPattern = "(load)( )(.*)"
	h.HandlerFunc = func(cmd *Command) string {
		arg := cmd.Args

		file, err := ioutil.ReadFile(arg)
		if err != nil {
			return err.Error()
		}
		fileText := string(file)
		loadedRequest.request = fileText
		return fileText
	}
	HandlerRegistry[h.CommandName] = h
}
