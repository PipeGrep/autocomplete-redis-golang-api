package main

import (
	"log"
	"time"
	"strings"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/fatih/structs"
	"github.com/garyburd/redigo/redis"
)

type request_body struct {
	Id      string
	Terms   []string
	Label   string
	Payload interface{}
	Score   float64
}

func IndexStore(res rest.ResponseWriter, req *rest.Request) {
	time_beginning := time.Now()
	contentType := req.PathParam("content-type")
	if len(contentType) == 0 {
		rest.Error(res, "Invalid content type", 400)
		return
	}

	payload := index__decode_request(res, req)
	if payload == nil {
		return
	}

	conn := redis_pool.Get()
	defer conn.Close()

	// Create token and/or add references to tokens in db
	nbr_tokens_references_added := 0
	for _, term := range payload.Terms {
		nbr_tokens_references_added += fill_token_index(conn, contentType, payload.Id, term, payload.Score)
	}

	if nbr_tokens_references_added > 0 {
		err := store_source(conn, contentType, payload.Id, payload.Label, payload.Payload)
		if err != nil {
			rest.Error(res, "Invalid document: "+err.Error(), 400)
			return
		}
	}

	log.Print("Indexed document '" + payload.Id + "', with terms: " + strings.Join(payload.Terms, ", "))
	res.WriteJson(map[string]interface{}{
		"acknowledge": true,
		"time": float64((time.Now().UnixNano() - time_beginning.UnixNano()) / 1000) / 1000,
		"number_tokens_indexed": nbr_tokens_references_added,
	})
}

func index__decode_request(res rest.ResponseWriter, req *rest.Request) *request_body {
	// set default values
	payload := request_body{
		Id:      "",
		Terms:   nil,
		Label:   "",
		Payload: nil,
		Score:   1,
	}

	err := req.DecodeJsonPayload(&payload)
	if err != nil {
		rest.Error(res, err.Error(), 400)
		return nil
	}

	if payload.Id == "" {
		rest.Error(res, "Bad format: Id", 400)
		return nil
	}
	if len(payload.Terms) == 0 {
		rest.Error(res, "Bad format: Terms is empty", 400)
		return nil
	}
	if payload.Label == "" {
		rest.Error(res, "Bad format: Label", 400)
		return nil
	}
	if payload.Score < 1 {
		rest.Error(res, "Bad format: Score", 400)
		return nil
	}

	for i, term := range payload.Terms {
		payload.Terms[i] = strings.ToLower(term)
	}

	payload.Label = strings.ToUpper(payload.Label[0:1]) + strings.ToLower(payload.Label[1:])

	return &payload
}

// call with -> content-type=users / id=4242 / term=google / score=3
// sorted set
// index:users:token:go  score+=3  values.push(id)
// index:users:token:goo  score+=3  values.push(id)
// index:users:token:goog  score+=3  values.push(id)
// index:users:token:googl  score+=3  values.push(id)
// index:users:token:google  score+=3  values.push(id)
func fill_token_index(conn redis.Conn, contentType string, id string, term string, score float64) int {
	ret := len(term) - config.Index.MinSizeToken
	for len(term) >= config.Index.MinSizeToken {
		key := "index:" + contentType + ":token:" + term
		conn.Do("ZADD", key, "INCR", score, id)
		term = term[0 : len(term)-1]
	}
	return ret
}

// call with -> content-type=users / id=4242 / label="Google is my friend" / payload={logo: "img/google.png", timestamp: 123456789}
// hash set
// index:users:source:4242  Label="Google is my friend"  Payload={logo: "img/google.png", timestamp: 123456789}
func store_source(conn redis.Conn, contentType string, id string, label string, payload interface{}) error {
	key := "index:" + contentType + ":source:" + id
	hash := Source{
		Id:      id,
		Label:   label,
		Payload: JsonToString(payload),
		Timestamp: time.Now().UnixNano(),
	}
	kvs := []interface{}{key}
	for key, value := range structs.Map(hash) {
		kvs = append(kvs, key, value)
	}
	_, err := conn.Do("HMSET", kvs...)
	return err
}
