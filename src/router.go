package main

import (
	"github.com/ant0ine/go-json-rest/rest"
)

func GetRouter() (rest.App, error) {
	return rest.MakeRouter(

		// GET /
		// reponse : { version: "1.0" }
		rest.Get("/", func(res rest.ResponseWriter, req *rest.Request) {
			res.WriteJson(map[string]string{"version": "1.0"})
		}),

		// PUT /index/#content-type/store
		//	{
		//		id: "4242",
		//		terms: ["simple", "example"],
		//		label: "It's just a simple example",
		//		payload: {
		//			foo: "bar"
		//		}
		//		score: 3		// OPTIONAL, default=1 (int)
		//	}
		//
		// Add a new entry in redis
		//	- id: can use few time the same ID to increment the score, but terms,
		//	  label and payload are stored only at first call
		rest.Put("/index/#content-type/store", IndexStore),

		// GET /index/#content-type/search?seed=examp&score=true&pretty=true&payload=true&ordered=true
		// -> Get a list of entries that match with the 'seed' parameter (multi-term matching)
		//	- ?seed= : terms to match with
		//	- ?ordered= (default true): return a list ordered by score
		//	- ?score= (default false): return score (cf /index/#content-type/store)
		//	- ?pretty= (default false): add <strong>exampl</strong> in the label
		//	- ?payload= (default false): return the payload in the response
		rest.Get("/index/#content-type/search", Search),

		// GET /index/#content-type/#id
		// -> Get an object by ID
		//	- ?payload= (default true): return the payload in the response
		rest.Get("/index/#content-type/#id", GetObject),
	)
}
