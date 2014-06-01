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

	url := ""
	if server.index == "" {
		url = fmt.Sprintf("http://%s:%s/%s", server.host, server.port, arg)
	} else {
		url = fmt.Sprintf("http://%s:%s/%s/%s", server.host, server.port, server.index, arg)
	}

	fmt.Println("Request:", url)
	res, err := getHttpResource(url)
	if err != nil {
		return err.Error()
	}
	return res
}

// curl -XPUT "http://localhost:9200/movies/movie/1" -d'{ ... body ... }''
// becomes
// put movie/1 { "title": "Alien", "director": "Ridley Scott", "year": 1979, "genres": ["Science fiction"] }
// Currently, must be on single line.
func handlePut(entry string) string {
	arg := strings.TrimPrefix(entry, CommandPut+" ")

	bodyIdx := strings.Index(arg, " ")
	queryArgs := arg[:bodyIdx]
	body := arg[bodyIdx:]

	url := fmt.Sprintf("http://%s:%s/%s/%s", server.host, server.port, server.index, queryArgs)
	fmt.Println("Request:", url)
	res, err := putHttpResource(url, body)
	if err != nil {
		return err.Error()
	}
	return res
}

// curl -XPOST "http://localhost:9200/movies/_search?pretty" -d'{ ... body ... }''
// becomes
// post _search?pretty { "query": { "term": { "director": "scott" } } }
// Currently, must be on single line.
func handlePost(entry string) string {
	arg := strings.TrimPrefix(entry, CommandPost+" ")

	bodyIdx := strings.Index(arg, " ")
	queryArgs := arg[:bodyIdx]
	body := arg[bodyIdx:]

	url := fmt.Sprintf("http://%s:%s/%s/%s", server.host, server.port, server.index, queryArgs)
	fmt.Println("Request:", url)
	res, err := postHttpResource(url, body)
	if err != nil {
		return err.Error()
	}
	return res
}
