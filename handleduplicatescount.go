package main

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"regexp"
	"sort"
)

// duplicatescount localhost:9200/srcindex/type/field
func init() {
	h := NewHandler()
	h.CommandName = "duplicatescount"
	h.CommandPattern = "(duplicatescount)(( )(.*))"
	h.HandlerFunc = func(cmd *Command) string {
		fmt.Println("Finding duplicates...")
		//args := strings.TrimPrefix(entry, CommandReindex+" ")
		args := cmd.Args

		// \w+|"[\w\s]*"
		r, err := regexp.Compile(`^(.*?):(\d+?)/(.*?)/(.*?)/(.*?)$`)
		if err != nil {
			return err.Error()
		}
		fmt.Println("Parsing command...")
		matches := r.FindAllStringSubmatch(args, -1)[0]
		fmt.Println("Parsed matches:", len(matches))
		srcHost := matches[1]
		srcPort := matches[2]
		srcIndex := matches[3]
		srcType := matches[4]
		srcField := matches[5]

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

		//dispMap(counts)
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
