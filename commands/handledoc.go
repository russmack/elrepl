package commands

import (
	"fmt"
	"github.com/russmack/elrepl/backend"
	"net/url"
	"regexp"
	"strings"
)

type DocCmd struct{}

// TODO: clean up.
func init() {
	h := backend.NewHandler()
	h.CommandName = "doc"
	h.CommandPattern = "(doc)(( )(.*))"
	h.Usage = "doc (get index docId) | (delete index type docId)"
	h.CommandParser = func(cmd *backend.Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Doc help
			regexp.MustCompile(`^doc /\?$`): func(s []string) (string, bool) {
				return "", false
			},
			// Get doc
			regexp.MustCompile(`^doc get ([a-zA-Z0-9\.\-_]+) ([a-zA-Z0-9\.\-_]+)$`): func(s []string) (string, bool) {
				fmt.Println("...", s[0])
				fmt.Println("...", s[1])
				fmt.Println("...", s[2])

				d := backend.Resource{
					Endpoint: "_mget",
					Scheme:   "http",
					Host:     backend.Server.Host,
					Port:     backend.Server.Port,
					Index:    s[1],
					Id:       s[2],
				}
				c := DocCmd{}
				r, ok := c.Get(d)
				return r, ok
			},
			// Delete doc
			regexp.MustCompile(`^doc delete ([a-zA-Z0-9\.\-_]+) ([a-zA-Z0-9\.\-_]+) ([a-zA-Z0-9\.\-_]+)$`): func(s []string) (string, bool) {
				d := backend.Resource{
					Endpoint: "_mget",
					Scheme:   "http",
					Host:     backend.Server.Host,
					Port:     backend.Server.Port,
					Index:    s[0],
					Type:     s[1],
					Id:       s[2],
				}
				c := DocCmd{}
				r, ok := c.Delete(d)
				return r, ok
			},
		}
		r, ok := h.Tokenize(strings.TrimSpace(cmd.Instruction), pattFn)
		return r, ok
	}
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		r, ok := h.CommandParser(cmd)
		if !ok {
			if r != "" {
				r += "\n\n"
			}
			return h.Usage, false
		}
		return r, true
	}
	backend.HandlerRegistry[h.CommandName] = h
}

func (c *DocCmd) Get(d backend.Resource) (string, bool) {
	postData := "{\"ids\": " + d.Id + " }"
	urlString := "post " + d.Host + " " + d.Port + " " + d.Index + "/" + d.Endpoint + " " + postData
	fmt.Println("...>>", urlString)

	cmdParser := backend.NewCommandParser()
	newCmd, err := cmdParser.Parse(urlString)
	if err != nil {
		return err.Error(), false
	}
	dispatcher := backend.NewDispatcher()
	res, ok := dispatcher.Dispatch(newCmd)
	return res, ok
}

//curl -XDELETE 'http://localhost:9200/twitter/tweet/1'
func (c *DocCmd) Delete(d backend.Resource) (string, bool) {
	u := new(url.URL)
	u.Scheme = d.Scheme
	u.Host = d.Host + ":" + d.Port
	u.Path = d.Index + "/" + d.Type + "/" + d.Id
	q := u.Query()
	q.Add("pretty", "true")
	u.RawQuery = q.Encode()

	//urlString = fmt.Sprintf("http://%s:%s/%s/%s/%s", backend.Server.Host, backend.Server.Port, indexName, typeName, id)
	fmt.Println("Request:", u)
	res, err := backend.DeleteHttpResource(u.String())
	if err != nil {
		return err.Error(), false
	}
	return res, true
}
