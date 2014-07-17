package main

import (
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
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

func handleUnknownEntry(cmd *Command) string {
	return fmt.Sprintf("Command not found: %s", cmd.Name)
}

func handleHostSet(cmd *Command) string {
	//arg := strings.TrimPrefix(entry, CommandServer+" ")
	arg := cmd.Args
	server.host = arg
	return "Set server host: " + arg
}

func handleHostGet() string {
	return "Server host: " + server.host
}

func handlePortSet(cmd *Command) string {
	//arg := strings.TrimPrefix(entry, CommandPort+" ")
	arg := cmd.Args
	server.port = arg
	return "Set server port: " + arg
}

func handlePortGet() string {
	return "Server port: " + server.port
}

func handleIndexSet(cmd *Command) string {
	//arg := strings.TrimPrefix(entry, CommandIndex+" ")
	arg := cmd.Args
	server.index = arg
	return "Set index: " + arg
}

func handleIndexGet() string {
	return "Index: " + server.index
}

func handleDir(cmd *Command) string {
	//arg := strings.TrimPrefix(entry, CommandDir+" ")
	arg := cmd.Args
	if arg == "" {
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

func handleLog(cmd *Command) string {
	logLevel = 1
	return "Logging level set to: " + strconv.Itoa(logLevel)
}

func handleLoad(cmd *Command) string {
	//arg := strings.TrimPrefix(entry, CommandLoad+" ")
	arg := cmd.Args

	file, err := ioutil.ReadFile(arg)
	if err != nil {
		return err.Error()
	}
	fileText := string(file)
	loadedRequest.request = fileText
	return fileText
}

func handleRun(cmd *Command) string {
	//arg := cmd.Args
	loadedParts := strings.SplitN(loadedRequest.request, "\n", 2)

	loadedCmdParts := strings.SplitN(loadedParts[0], " ", 2)
	loadedCmd := loadedCmdParts[0]
	//loadedArgs := loadedCmdParts[1]

	//loadedQuery := loadedParts[1]

	cmdParser := NewCommandParser()
	newCmd, err := cmdParser.Parse(loadedRequest.request)
	if err != nil {
		return "Unable to parse loaded query for run command."
	}
	if strings.ToLower(loadedCmd) == "post" {
		resp := handlePost(newCmd)
		return resp
	} else if strings.ToLower(loadedCmd) == "put" {
		resp := handlePut(newCmd)
		return resp
	} else if strings.ToLower(loadedCmd) == "get" {
		resp := handleGet(newCmd)
		return resp
	}
	return "Unable to run loaded query."
}

func handleGet(cmd *Command) string {
	//arg := strings.TrimPrefix(entry, CommandGet+" ")
	arg := cmd.Args

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
func handlePut(cmd *Command) string {
	//arg := strings.TrimPrefix(entry, CommandPut+" ")
	arg := cmd.Args

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

// http://localhost:9200/movies/_search
// localhost:9200/movies/_search
// /movies/_search
// movies/_search
// http://localhost:9200/movies/movie/_search
// localhost:9200/movies/movie/_search
// /movies/movie/_search
// movies/movie/_search
// movie/_search

// http://localhost:9200/movies/_search
// http://localhost:9200/movies/movie/_search
// localhost:9200/movies/movie/_search
// localhost:9200/movies/_search
// movies/movie/_search
// movies/_search
// movie/_search

// curl -XPOST "http://localhost:9200/movies/_search?pretty" -d'{ ... body ... }''
// becomes
// post _search?pretty { "query": { "term": { "director": "scott" } } }
// Currently, must be on single line.
func handlePost(cmd *Command) string {
	queryHost := server.host
	queryPort := server.port
	//queryIndex := server.index

	// TODO: being lazy, for now.  Implement parsing to get this from the query.
	// This url parsing and construction is a bit of a mess.
	if queryHost == "" || queryPort == "" {
		return "Please set the host and port using 'host xxxxx' and 'port nnnn'."
	}

	if strings.HasPrefix(cmd.Args, "http://") || strings.HasPrefix(cmd.Args, "https://") {
		return "Full url parsing not yet implemented. Please remove 'http://host:port' prefix from url."
	}

	//arg := strings.TrimPrefix(entry, CommandPost+" ")
	arg := cmd.Args
	arg = strings.Replace(arg, "\r\n", " ", -1)
	arg = strings.Replace(arg, "\n", " ", -1)
	bodyIdx := strings.Index(arg, "{")
	queryArgs := arg[:bodyIdx]
	body := arg[bodyIdx:]
	//body = strings.Replace(body, "\n", "", -1)

	queryArgs = strings.TrimPrefix(arg, "/")
	/*
		if queryIndex == "" {
			cleanArg := strings.TrimPrefix(arg, "/")
			if cleanArg != "" {
				argParts := strings.Split(cleanArg, "/")
				if len(argParts) > 0 {
					queryIndex = argParts[0]
				}
			}
		}
	*/
	fmt.Println("queryHost:", queryHost)
	fmt.Println("queryPort:", queryPort)
	//fmt.Println("queryIndex:", queryIndex)
	fmt.Println("queryArgs:", queryArgs)
	//url := fmt.Sprintf("http://%s:%s/%s/%s", queryHost, queryPort, queryIndex, queryArgs)
	url := fmt.Sprintf("http://%s:%s/%s", queryHost, queryPort, queryArgs)
	fmt.Println("Request:", url)
	fmt.Println("----")
	res, err := postHttpResource(url, body)
	if err != nil {
		return err.Error()
	}
	return res
}

// reindex localhost:9200/srcindex/type localhost:9200/targetindex/routing
func handleReindex(cmd *Command) string {
	//args := strings.TrimPrefix(entry, CommandReindex+" ")
	args := cmd.Args

	// \w+|"[\w\s]*"
	r, err := regexp.Compile(`^(.*?):(\d+?)/(.*?)/(.*?)/? (.*?):(\d+?)/(.*?)(/(.*?))?$`)
	if err != nil {
		return err.Error()
	}
	matches := r.FindAllStringSubmatch(args, -1)[0]
	srcHost := matches[1]
	srcPort := matches[2]
	srcIndex := matches[3]
	srcType := matches[4]
	tgtHost := matches[5]
	tgtPort := matches[6]
	tgtIndex := matches[7]
	tgtRouting := matches[8]

	api.Domain = srcHost
	api.Port = srcPort

	scanArgs := map[string]interface{}{"search_type": "scan", "scroll": "1m", "size": "1000"}
	scanResult, err := core.SearchUri(srcIndex, srcType, scanArgs)
	if err != nil {
		fmt.Println("Failed getting scan result for index:", srcIndex, "; err:", err)
		return err.Error()
	}

	//total := scanResult.Hits.Total

	scrollId := scanResult.ScrollId
	counter := 0
	failures := 0

	scrollArgs := map[string]interface{}{"scroll": "1m"}
	scrollResult, err := core.Scroll(scrollArgs, scrollId)
	if err != nil {
		fmt.Println("Failed getting scroll result for index:", srcIndex, "; err:", err)
		return err.Error()
	}

	var indexArgs map[string]interface{} = nil
	if tgtRouting != "" {
		indexArgs = map[string]interface{}{"routing": tgtRouting}
	}
	for len(scrollResult.Hits.Hits) > 0 {
		fmt.Println("Scroll result hits:", len(scrollResult.Hits.Hits))
		for _, j := range scrollResult.Hits.Hits {
			api.Domain = tgtHost
			api.Port = tgtPort
			_, err := core.Index(tgtIndex, srcType, j.Id, indexArgs, j.Source)
			if err != nil {
				fmt.Println("Failed inserting document, id:", j.Id, "; ", err)
				failures++
				continue
			}
			counter++
		}

		api.Domain = srcHost
		api.Port = srcPort
		// ScrollId changes with every request.
		scrollId = scrollResult.ScrollId
		scrollArgs := map[string]interface{}{"scroll": "1m"}
		scrollResult, err = core.Scroll(scrollArgs, scrollId)
		if err != nil {
			fmt.Println("Failed getting scroll result for index:", srcIndex, "; err:", err)
			return err.Error()
		}
	}
	return fmt.Sprintf("Total processed: %d.  %d failed.", counter, failures)
}
