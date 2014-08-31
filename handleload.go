package main

import (
	"io/ioutil"
)

func init() {
	h := NewHandler()
	h.CommandName = "load"
	h.CommandPattern = "(load)( )(.*)"
	h.Usage = "load filename"
	h.HandlerFunc = func(cmd *Command) string {
		filename := cmd.Tokens[1]
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			return err.Error()
		}
		fileText := string(file)
		loadedRequest.request = fileText
		return fileText
	}
	HandlerRegistry[h.CommandName] = h
}
