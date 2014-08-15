package main

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"sort"
)

func init() {
	h := NewHandler()
	h.CommandName = "duplicatescount"
	h.CommandPattern = "(duplicatescount)(( )(.*))"
	h.Usage = "duplicatecount [host port] index type field"
	h.CommandParser = func(cmd *Command) (ParseMap, bool) {
		p := ParseMap{}
		switch len(cmd.Args) {
		case 3:
			p["host"] = server.host
			p["port"] = server.port
			p["index"] = cmd.Args[0]
			p["type"] = cmd.Args[1]
			p["field"] = cmd.Args[2]
		case 5:
			p["host"] = cmd.Args[0]
			p["port"] = cmd.Args[1]
			p["index"] = cmd.Args[2]
			p["type"] = cmd.Args[3]
			p["field"] = cmd.Args[4]
		default:
			//case "/?"
			//case ""
			return p, false
		}
		return p, true
	}
	h.HandlerFunc = func(cmd *Command) string {
		fmt.Println("Finding duplicates...")
		p, ok := h.CommandParser(cmd)
		if !ok {
			return usageMessage(h.Usage)
		}
		srcHost := p["host"]
		srcPort := p["port"]
		srcIndex := p["index"]
		srcType := p["tyep"]
		srcField := p["field"]

		api.Domain = srcHost
		api.Port = srcPort

		fmt.Println("Scanning...")
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

		fmt.Println("Scrolling...")
		scrollArgs := map[string]interface{}{"scroll": "1m"}
		scrollResult, err := core.Scroll(scrollArgs, scrollId)
		if err != nil {
			fmt.Println("Failed getting scroll result for index:", srcIndex, "; err:", err)
			return err.Error()
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
						if k == srcField {
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
				fmt.Println("Failed getting scroll result for index:", srcIndex, "; err:", err)
				return err.Error()
			}
		}
		dispPairList(counts)
		return fmt.Sprintf("Total processed: %d.  %d failed.", counter, failures)
	}
	HandlerRegistry[h.CommandName] = h
}

type Duplicate struct {
	IdList []string
	Value  string
	Count  int
}

func dispMap(counts map[string]int) {
	for k, v := range counts {
		fmt.Print(k, ":", v, "|")
	}
}

func dispPairList(counts map[string]*Duplicate) {
	pairlist := sortMapByValue(counts)
	tot := len(pairlist)
	fmt.Println("total:", tot)
	duplicated := 0
	for i, _ := range pairlist {
		if pairlist[i].Duplicate.Count > 1 {
			duplicated++
		}
	}
	for i := tot - 1; i > tot-30; i-- {
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
