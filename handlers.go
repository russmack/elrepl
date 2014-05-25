package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func handleHelp() string {
	return `
	Help
	----
	Commands:
	  eg:
	  host localhost
	  port 9200
	  index movies
	  get _search?q=title:thx1138
	`
}

func handleVersion() string {
	return `
	elRepl version 0.1
	`
}

func handleExit() string {
	fmt.Println("Bye.")
	os.Exit(0)
	return ""
}

func handleUnknownEntry(entry string) string {
	return fmt.Sprintf("Command not found: %s", entry)
}

func handleServerSet(entry string) string {
	arg := strings.TrimPrefix(entry, CommandServer+" ")
	server.host = arg
	return "Set server host: " + arg
}

func handleServerGet() string {
	return "Server host: " + server.host
}

func handlePortSet(entry string) string {
	arg := strings.TrimPrefix(entry, CommandPort+" ")
	server.port = arg
	return "Set server port: " + arg
}

func handlePortGet() string {
	return "Server port: " + server.port
}

func handleIndexSet(entry string) string {
	arg := strings.TrimPrefix(entry, CommandIndex+" ")
	server.index = arg
	return "Set index: " + arg
}

func handleIndexGet() string {
	return "Index: " + server.index
}

func handleDir(entry string) string {
	arg := strings.TrimPrefix(entry, CommandDir+" ")
	if arg == "dir" {
		arg = "."
	}
	dirFiles, err := ioutil.ReadDir(arg)
	if err != nil {
		return err.Error()
	}
	files := ""
	for _, j := range dirFiles {
		files += j.Name() + "\n"
	}
	return files
}

func handleLoad(entry string) string {
	arg := strings.TrimPrefix(entry, CommandLoad+" ")
	file, err := ioutil.ReadFile(arg)
	if err != nil {
		return err.Error()
	}
	fileText := string(file)
	loadedRequest.request = fileText
	return fileText
}

func handleGet(entry string) string {
	arg := strings.TrimPrefix(entry, CommandGet+" ")
	//return "Not implemented, cannot execute " + entry
	url := fmt.Sprintf("http://%s:%s/%s/%s", server.host, server.port, server.index, arg)
	fmt.Println("Request:", url)
	res, err := getHttpResource(url)
	if err != nil {
		return err.Error()
	}
	return res
}
