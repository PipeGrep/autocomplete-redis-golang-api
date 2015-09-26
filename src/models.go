package main

import (
	"sort"
)

type Source struct {
	Id      string
	Label   string
	Payload string
	Timestamp int64
}

type Index struct {
	Id    string
	Score float64
}

type SortSourceByScore struct {
	Keys   []string
	Values map[string]map[string]interface{}
}

func (s SortSourceByScore) Len() int {
	return len(s.Keys)
}
func (s SortSourceByScore) Swap(i, j int) {
	s.Values[s.Keys[i]], s.Values[s.Keys[j]] = s.Values[s.Keys[j]], s.Values[s.Keys[i]]
	s.Keys[i], s.Keys[j] = s.Keys[j], s.Keys[i]
}
func (s SortSourceByScore) Less(i, j int) bool {
	return s.Values[s.Keys[i]]["Score"].(float64) > s.Values[s.Keys[j]]["Score"].(float64)
}

func SortSourceByScoreRev(list map[string]map[string]interface{}) map[string]map[string]interface{} {
	sort_list_result := new(SortSourceByScore)
	sort_list_result.Values = list
	for key, _ := range list {
		sort_list_result.Keys = append(sort_list_result.Keys, key)
	}
	sort.Sort(sort_list_result)
	sort.Reverse(sort_list_result)
	return sort_list_result.Values
}
