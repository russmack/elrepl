package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Server struct {
	host  string
	port  string
	index string
}

type LoadedRequest struct {
	request string
}

const (
	CommandServer = "host"
	CommandPort   = "port"
	CommandIndex  = "index"
	CommandDir    = "dir"
	CommandLoad   = "load"
	CommandGet    = "get"
	CommandPost   = "post"
)

var (
	//Stdin *bufio.Reader
	server        = Server{}
	loadedRequest = LoadedRequest{}
)

func init() {
	server.host = "localhost"
	server.port = "9200"
}

func main() {
	displayWelcome()
	reploop()
}

func displayWelcome() {
	message := `
	el Repl
	=======

	Welcome to elRepl, an elasticsearch repl.
	`

	fmt.Println(message)
}

func reploop() {
	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		entered, err := in.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		//entry := strings.TrimLeft(entered[:len(entered)-1], "\t ")
		entry := strings.Trim(entered, "\t \r\n")
		output := dispatch(entry)
		if len(output) > 0 {
			fmt.Println(output)
		}
	}

}

func dispatch(entry string) string {
	if len(entry) == 0 {
		return entry
	}
	switch {
	case entry == "?" || entry == "help":
		return handleHelp()
	case entry == "version":
		return handleVersion()
	case entry == "exit" || entry == "quit" || entry == "bye":
		return handleExit()
	case strings.HasPrefix(entry, CommandServer+" "):
		return handleServerSet(entry)
	case entry == CommandServer:
		return handleServerGet()
	case strings.HasPrefix(entry, CommandPort+" "):
		return handlePortSet(entry)
	case entry == CommandPort:
		return handlePortGet()
	case strings.HasPrefix(entry, CommandIndex+" "):
		return handleIndexSet(entry)
	case entry == CommandIndex:
		return handleIndexGet()
	case strings.HasPrefix(entry, CommandDir):
		return handleDir(entry)
	case strings.HasPrefix(entry, CommandLoad):
		return handleLoad(entry)
	case strings.HasPrefix(entry, CommandGet):
		return handleGet(entry)
	default:
		return handleUnknownEntry(entry)
	}
}

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

func getHttpResource(url string) (string, error) {
	client := NewTimeoutClient(10*time.Second, 15*time.Second)

	ok := false
	retriesAllowed := 3
	resp := &http.Response{}
	var err error = nil

	for i := 0; i < retriesAllowed; i++ {
		//fmt.Println("Attempt", i+1, "of", retriesAllowed)

		resp, err = client.Get(url)
		if err != nil {
			fmt.Println("err, getHttpResource, get: ", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			ok = true
			break
		} else {
			errmsg := fmt.Sprintf("Http server returned status code:", resp.StatusCode)
			err = errors.New(errmsg)
		}
	}

	body := []byte{}
	if ok {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("err, getHttpResource, readall: ", err)
		}
	}

	//return bytes.NewReader(body)
	return string(body), err
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		err = conn.SetDeadline(time.Now().Add(rwTimeout))
		if err != nil {
			fmt.Println("Unable to set deadline of http connection:", err)
			return nil, err
		}
		return conn, nil
	}
}

func NewTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
		},
	}
}
