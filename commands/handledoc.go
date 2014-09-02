package commands

import (
	"fmt"
	"github.com/russmack/elrepl/backend"
	"net/url"
	"strings"
)

type DocCmd struct{}

// TODO: clean up.
func init() {
	h := backend.NewHandler()
	h.CommandName = "doc"
	h.CommandPattern = "(doc)(( )(.*))"
	h.Usage = "doc (get index docId) | (delete index type docId)"
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		//argParts := strings.SplitN(cmd.Args, " ", 3)
		argParts := cmd.Tokens

		arg := "_mget?pretty"
		urlString := ""
		res := ""
		ok := false
		fmt.Println("ARGS:", argParts)
		if argParts[0] == "get" {
			indexName := argParts[1]
			ids := argParts[2]
			//urlString = fmt.Sprintf("http://%s:%s", backend.Server.Host, backend.Server.Port)
			//indexName := argParts[1]
			//aliasName := argParts[2]
			// "ids" : ["1", "1"]
			urlString = "post " + indexName + "/" + arg + " " + "{\"ids\": " + ids + " }"
			cmdParser := backend.NewCommandParser()
			newCmd, err := cmdParser.Parse(urlString)
			if err != nil {
				return err.Error(), false
			}
			dispatcher := backend.NewDispatcher()
			res, ok = dispatcher.Dispatch(newCmd)
		} else if argParts[0] == "delete" {
			//curl -XDELETE 'http://localhost:9200/twitter/tweet/1'
			indexName := argParts[1]
			delArgs := strings.Split(argParts[2], " ")
			typeName := delArgs[0]
			id := delArgs[1]

			u := new(url.URL)
			u.Scheme = "http"
			u.Host = backend.Server.Host + ":" + backend.Server.Port
			u.Path = indexName + "/" + typeName + "/" + id
			q := u.Query()
			q.Add("pretty", "true")
			u.RawQuery = q.Encode()

			//urlString = fmt.Sprintf("http://%s:%s/%s/%s/%s", backend.Server.Host, backend.Server.Port, indexName, typeName, id)
			fmt.Println("Request:", u)
			delRes, err := backend.DeleteHttpResource(u.String())
			if err != nil {
				return err.Error(), false
			}
			res = delRes
			ok = true
		}
		//fmt.Println("Request:", url)
		//res, err :=backend.GetHttpResource(url)
		//if err != nil {
		//	return err.Error(), false
		//}
		return res, ok
	}
	backend.HandlerRegistry[h.CommandName] = h
}
