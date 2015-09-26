package main

import "strings"
import "encoding/json"

type SortByLengthRev []string

func (s SortByLengthRev) Len() int {
	return len(s)
}
func (s SortByLengthRev) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortByLengthRev) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func IndexMulti(s string, token string) []int {
	ret := []int{}

	sum := 0
	for strings.Index(s, token) != -1 {
		i := strings.Index(s, token)
		ret = append(ret, i+sum)
		s = s[i+1:]
		sum += i + 1
	}

	return ret
}

func JsonToString(j interface{}) string {
	str, err := json.Marshal(j)

	if err != nil {
		return ""
	}
	return string(str[0:len(str)])
}

func StringToJson(str string) interface{} {
	var j interface{}
	err := json.Unmarshal([]byte(str), &j)

	if err != nil {
		return nil
	}
	return j
}
