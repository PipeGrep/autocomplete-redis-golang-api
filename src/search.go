package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/garyburd/redigo/redis"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Search(res rest.ResponseWriter, req *rest.Request) {
	time_beginning := time.Now()
	contentType, seeds, flags, err := search__decode_request(res, req)
	if err == true {
		return
	}

	conn := redis_pool.Get()
	defer conn.Close()

	set_list_results := make(map[string]map[string]interface{})

	// search seeds in the index
	for _, seed := range seeds {
		key := "index:" + contentType + ":token:" + seed
		values, err := redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "WITHSCORES"))
		if err != nil {
			rest.Error(res, "Unexpected error: " + err.Error(), 400)
			return
		}

		// in the values array:
		// 	- even items are ids (a reference to a document)
		//	- odd items are their respective score
		i := 0
		for i < len(values) {
			doc_id := values[i]
			score, _ := strconv.ParseFloat(values[i + 1], 64)

			source, err := get_source_by_id(res, conn, contentType, doc_id)
			if err == true {
				return
			}
			if source == nil {
				continue
			}

			source["Score"] = score

			// flag payload
			if flags["payload"] == false {
				source["Payload"] = nil
			}

			if set_list_results[doc_id] == nil {
				// flag pretty
				if flags["pretty"] == true {
					source["LabelPretty"] = get_pretty_label(source["Label"].(string), seeds, "<strong>", "</strong>")
				}
				set_list_results[doc_id] = source
			} else {
				set_list_results[doc_id]["Score"] = set_list_results[doc_id]["Score"].(float64) + source["Score"].(float64)
			}
			i += 2
		}
	}

	set_list_results = SortSourceByScoreRev(set_list_results)

	// flag score
	if flags["Score"] == false {
		for key, _ := range set_list_results {
			set_list_results[key]["Score"] = nil
		}
	}

	// without IDs as key
	// return an array, not an object
	list_results := []map[string]interface{}{}
	for _, result := range set_list_results {
		list_results = append(list_results, result)
	}

	// response
	res.WriteJson(map[string]interface{}{
		"request": map[string]interface{}{
			"content-type": contentType,
			"seeds": seeds,
			"flags": flags,
		},
		"time": float64((time.Now().UnixNano() - time_beginning.UnixNano()) / 1000) / 1000,
		"number_results": len(list_results),
		"results": list_results,
	})
}

func search__decode_request(res rest.ResponseWriter, req *rest.Request) (string, []string, map[string]bool, bool) {
	// get CONTENT-TYPE parameter
	contentType := req.PathParam("content-type")
	if len(contentType) == 0 {
		rest.Error(res, "Invalid content type", 400)
		return "", nil, nil, true
	}

	// get SEED parameter
	if req.FormValue("seed") == "" {
		rest.Error(res, "Add the beginning of the keyword in the query ?seed=", 400)
		return "", nil, nil, true
	}
	seeds := strings.Split(strings.ToLower(req.FormValue("seed")), " ")

	// get ORDERED, SCORE, PRETTY and PAYLOAD parameters
	flags := map[string]bool{
		"ordered": true,
		"score": false,
		"pretty": false,
		"payload": false,
	}
	for flag_name, _ := range flags {
		if req.Form[flag_name] != nil {
			flag_str := req.FormValue(flag_name)
			flag, err := strconv.ParseBool(flag_str)
			if err != nil {
				rest.Error(res, "Invalid query " + flag_name + ". " + err.Error(), 400)
				return "", nil, nil, true
			}
			flags[flag_name] = flag
		}
	}

	return contentType, seeds, flags, false
}

// add some bold html tag around seed in label
// The with an upper case as first character
func get_pretty_label(label string, seeds []string, opening_tag string, closing_tag string) string {
	pretty_label := strings.ToLower(label)

	// seeds must be sorted by size (large to small), because "goo" is included in "google"
	sort.Sort(SortByLengthRev(seeds))

	for _, seed := range seeds {
		pretty_label = strings.Replace(pretty_label, seed, opening_tag+seed+closing_tag, -1)

		// in case of seed equal to "t" -> <s<strong>t</strong>ong>
		if strings.Index(opening_tag, seed) != -1 {
			seed_in_opening_tag := strings.Replace(opening_tag, seed, opening_tag+seed+closing_tag, -1)
			pretty_label = strings.Replace(pretty_label, seed_in_opening_tag, opening_tag, -1)
		}
		if strings.Index(opening_tag, seed) != -1 {
			seed_in_closing_tag := strings.Replace(closing_tag, seed, opening_tag+seed+closing_tag, -1)
			pretty_label = strings.Replace(pretty_label, seed_in_closing_tag, closing_tag, -1)
		}
	}

	// delete multiple inclusion of html tags (ex: <strong><strong>goo</strong>gle</strong>
	pretty_label, err := epur_tag_inclusions(pretty_label, opening_tag, closing_tag)
	if err == true {
		return label
	}

	// we cannot up the first character of pretty_label because it can be a '<'
	// we cannot use the first character of label to find the first character of pretty_label because 'str' is include in <strong>
	if len(pretty_label) > 3 && strings.Index(pretty_label, opening_tag) == 0 {
		first_character_label_upper_case := strings.ToUpper(string(pretty_label[len(opening_tag)]))
		pretty_label = opening_tag + first_character_label_upper_case + pretty_label[len(opening_tag)+1:]
	} else {
		first_character_label_upper_case := strings.ToUpper(string(pretty_label[0]))
		pretty_label = first_character_label_upper_case + pretty_label[1:]
	}

	return pretty_label
}

// delete multiple inclusion of html tags (ex: <strong><strong>goo</strong>gle</strong>
func epur_tag_inclusions(pretty_label string, opening_tag string, closing_tag string) (string, bool) {
	opening_tags_indexes := IndexMulti(pretty_label, opening_tag)
	closing_tags_indexes := IndexMulti(pretty_label, closing_tag)
	if len(opening_tags_indexes) != len(closing_tags_indexes) {
		return "", true
	}
	i := 0
	for i < len(opening_tags_indexes)-1 {
		if opening_tags_indexes[i+1] < closing_tags_indexes[i] {
			// delete closing tag
			pretty_label = pretty_label[:closing_tags_indexes[i]] + pretty_label[closing_tags_indexes[i]+len(closing_tag):]
			// delete opening tag
			pretty_label = pretty_label[:opening_tags_indexes[i+1]] + pretty_label[opening_tags_indexes[i+1]+len(opening_tag):]
			// reset the loop
			opening_tags_indexes = IndexMulti(pretty_label, opening_tag)
			closing_tags_indexes = IndexMulti(pretty_label, closing_tag)
			i = 0
		} else {
			i++
		}
	}
	return pretty_label, false
}
