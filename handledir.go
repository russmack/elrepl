package main

import (
	"io/ioutil"
)

func init() {
	h := NewHandler()
	h.CommandName = "dir"
	h.CommandPattern = "(dir)( )(.*)"
	h.Usage = "dir"
	h.HandlerFunc = func(cmd *Command) string {
		dir := ""
		if len(cmd.Tokens) == 1 {
			dir = "."
		} else {
			dir = cmd.Tokens[1]
		}
		dirFiles, err := ioutil.ReadDir(dir)
		if err != nil {
			return err.Error()
		}
		files := ""
		for _, j := range dirFiles {
			files += j.Name() + "\n"
		}
		return files
	}
	HandlerRegistry[h.CommandName] = h
}
