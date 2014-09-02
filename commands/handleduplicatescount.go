package commands

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"github.com/russmack/elrepl/backend"
	"regexp"
	"sort"
	"strings"
)

type DuplicatesCountCmd struct{}

func init() {
	h := backend.NewHandler()
	h.CommandName = "duplicatescount"
	h.CommandPattern = "(duplicatescount)(( )(.*))"
	h.Usage = "duplicatescount host port index type field"
	h.CommandParser = func(cmd *backend.Command) (string, bool) {
		pattFn := map[*regexp.Regexp]func([]string) (string, bool){
			// Duplicatescount help
			regexp.MustCompile(`^duplicatescount /\?$`): func(s []string) (string, bool) {
				return "", false
			},
			// Count duplicates
			regexp.MustCompile(`^duplicatescount ([a-zA-Z0-9\.\-]+) ([0-9]{1,5}) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+) ([a-zA-Z0-9\.\-]+)$`): func(s []string) (string, bool) {
				d := backend.Resource{
					Host:  s[1],
					Port:  s[2],
					Index: s[3],
					Type:  s[4],
					Field: s[5],
				}
				c := DuplicatesCountCmd{}
				r, ok := c.Do(d)
				return r, ok
			},
		}
		r, ok := h.Tokenize(strings.TrimSpace(cmd.Instruction), pattFn)
		return r, ok
	}
	h.HandlerFunc = func(cmd *backend.Command) (string, bool) {
		fmt.Println("Finding duplicates...")
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

func (c *DuplicatesCountCmd) Do(d backend.Resource) (string, bool) {
	api.Domain = d.Host
	api.Port = d.Port
	fmt.Println("Scanning...")
	scanArgs := map[string]interface{}{"search_type": "scan", "scroll": "1m", "size": "1000"}
	scanResult, err := core.SearchUri(d.Index, d.Type, scanArgs)
	if err != nil {
		fmt.Println("Failed getting scan result for index:", d.Index, "; err:", err)
		return err.Error(), false
	}
	//total := scanResult.Hits.Total
	scrollId := scanResult.ScrollId
	counter := 0
	failures := 0
	fmt.Println("Scrolling...")
	scrollArgs := map[string]interface{}{"scroll": "1m"}
	scrollResult, err := core.Scroll(scrollArgs, scrollId)
	if err != nil {
		fmt.Println("Failed getting scroll result for index:", d.Index, "; err:", err)
		return err.Error(), false
	}
	fmt.Println("Counting...")
	counts := make(map[string]*Duplicate)
	var f interface{}
	for len(scrollResult.Hits.Hits) > 0 {
		fmt.Println("Scroll result hits:", len(scrollResult.Hits.Hits))
		for _, j := range scrollResult.Hits.Hits {
			docId := j.Id
			err := json.Unmarshal(*j.Source, &f)
			if err != nil {
				fmt.Println("ERR:", err)
			}
			m := f.(map[string]interface{})
			for k, v := range m {
				switch vv := v.(type) {
				case string:
					if k == d.Field {
						_, ok := counts[vv]
						if ok {
							counts[vv].Count++
							counts[vv].IdList = append(counts[vv].IdList, docId)
						} else {
							counts[vv] = &Duplicate{}
							counts[vv].Count = 1
							counts[vv].Value = vv
							counts[vv].IdList = []string{docId}
						}
					}
				default:
					//nothing
				}
			}
			counter++
		}
		// ScrollId changes with every request.
		scrollId = scrollResult.ScrollId
		scrollArgs := map[string]interface{}{"scroll": "1m"}
		scrollResult, err = core.Scroll(scrollArgs, scrollId)
		if err != nil {
			fmt.Println("Failed getting scroll result for index:", d.Index, "; err:", err)
			return err.Error(), false
		}
	}
	showPairList(counts)
	return fmt.Sprintf("Total processed: %d.  %d failed.", counter, failures), true
}

type Duplicate struct {
	IdList []string
	Value  string
	Count  int
}

func showMap(counts map[string]int) {
	for k, v := range counts {
		fmt.Print(k, ":", v, "|")
	}
}

func showPairList(counts map[string]*Duplicate) {
	pairlist := sortMapByValue(counts)
	tot := len(pairlist)
	fmt.Println("total:", tot)
	duplicated := 0
	for i, _ := range pairlist {
		if pairlist[i].Duplicate.Count > 1 {
			duplicated++
		}
	}
	displayMax := 30
	if tot < displayMax {
		displayMax = tot
	}
	fmt.Println("Displaying max:", displayMax)
	for i := tot - 1; i >= tot-displayMax; i-- {
		if pairlist[i].Duplicate.Count > 1 {
			fmt.Println(pairlist[i].Duplicate.Count, " : ", pairlist[i].Key, " : ", pairlist[i].Duplicate.IdList)
		}
	}
	fmt.Println(duplicated, "duplicated.")
}

type Pair struct {
	Key       string
	Duplicate *Duplicate
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Duplicate.Count < p[j].Duplicate.Count }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]*Duplicate) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}
