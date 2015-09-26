package main

import (
        "github.com/ant0ine/go-json-rest/rest"
	"github.com/garyburd/redigo/redis"
	"github.com/fatih/structs"
	"time"
	"strconv"
)

func GetObject(res rest.ResponseWriter, req *rest.Request) {
	time_beginning := time.Now()
	contentType, doc_id, flags, err := get_object__decode_request(res, req)
	if err == true {
		return
	}

	conn := redis_pool.Get()
	defer conn.Close()

	source, err := get_source_by_id(res, conn, contentType, doc_id)
	if err == true {
		return
	}

	// flag payload
	if flags["payload"] == false {
		source["Payload"] = nil
	}

	res.WriteJson(map[string]interface{}{
		"request": map[string]interface{}{
			"content-type": contentType,
			"flags": flags,
		},
                "time": float64((time.Now().UnixNano() - time_beginning.UnixNano()) / 1000) / 1000,
		"doc": source,
	})
}

func get_object__decode_request(res rest.ResponseWriter, req *rest.Request) (string, string, map[string]bool, bool) {
	// get CONTENT-TYPE parameter
	contentType := req.PathParam("content-type")
	if len(contentType) == 0 {
		rest.Error(res, "Invalid content type", 400)
		return "", "", nil, true
	}

	// get ID parameter
	doc_id := req.PathParam("id")
	if len(doc_id) == 0 {
		rest.Error(res, "Invalid content ID", 400)
		return "", "", nil, true
	}

	// get PAYLOAD parameters
	flags := map[string]bool{
		"payload": true,
	}
	for flag_name, _ := range flags {
		if req.Form[flag_name] != nil {
			flag_str := req.FormValue(flag_name)
			flag, err := strconv.ParseBool(flag_str)
			if err != nil {
				rest.Error(res, "Invalid query " + flag_name + ". " + err.Error(), 400)
				return "", "", nil, true
			}
			flags[flag_name] = flag
		}
	}

	return contentType, doc_id, flags, false
}

func get_source_by_id(res rest.ResponseWriter, conn redis.Conn, contentType string, doc_id string) (map[string]interface{}, bool) {
	hash_key := "index:" + contentType + ":source:" + doc_id
	hash_value := Source{"", "", "", 0}

	exists, err := redis.Bool(conn.Do("EXISTS", hash_key))
	if err != nil {
		rest.Error(res, "Unexpected error. " + err.Error(), 400)
		return nil, true
	}
	if exists == false {
		return nil, false
	}

	// get the document indexed
	values, err := conn.Do("HGETALL", hash_key)
	if err == nil {
		err = redis.ScanStruct(values.([]interface{}), &hash_value)
	}
	if err != nil {
		rest.Error(res, "Unexpected error. " + err.Error(), 400)
		return nil, true
	}

	source := structs.Map(hash_value)

	source["Id"] = doc_id
	source["Payload"] = StringToJson(source["Payload"].(string))

	return source, false
}
