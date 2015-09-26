Search as you Type API (Golang + Redis)
============================================

# Features

## Index
- Separate **content-types** in different indexes (ex: usernames, article titles...)
- Can index **multiple terms** in the same string
- Attach a **payload**/document to the indexed string
- **Scoring mechanism** :
  - by indexing many times the same document
  - by indexing a score
- Indexation protected by a basic authentication

## Search
- **Sorted** results (by score)
- **Case insensitive**
- Pretty mode (add a <strong></strong> tag around the terms)
- Return a payload
- Go + Redis == **Very fast** (Elasticsearch run out of time)

## Other names :
- "Incremental search API"
- "Autocompletion API"
- "As you type suggestions API"
- "Instant results as you type API"

# Install

## Configuration file

Many values can be set by the environment :
- redis.hostname
- redis.port
- redis.port
- http.port
- http.basic.username
- http.basic.password

Example:
"""js
{
    "redis": {
        "hostname": "$REDIS_HOSTNAME",
        "port": "$REDIS_PORT",
        "password": "$REDIS_PASSWORD"
    },
    "http": {
        "port": "$API_PORT",
        "basic": {
            "username": "$API_BASIC_USERNAME",
            "password": "$API_BASIC_PASSWORD"
        }
    },
    "index": {
                "minSizeToken": 2,
                "prefix": "users"
    }
}
"""


## With docker

Build image :
"""sh
docker build -t autocompletion_api .
"""

Run container :
"""sh
$ docker run -d -p 6379:6379 --name redis redis
$ docker run -d -p 8080:8080 --link redis:redis autocompletion_api
"""

## Without docker

Update config.json setting redis hostname to "localhost"

"""sh
$ make vendor_get
$ make run
"""

# HOW TO

## Build the index

PUT /index/#content-type/store

- **content-type**: "username", "article-title" or whatever...
- request **body** (json):
  - **id** (string): ID of the document to index. If used multiple time, the score increase of 1.
  - **terms** ([]string): You can index a string with many world. Example, my twitter can be indexed with ["@SamuelBerthe", "Samuel", Berthe"].
  - **payload** (json, optional): whatever you want to store with the string.
  - **score** (float, optional, default=1): the search result can be sorted by score

### Return

"""js
{ results: [...] }
"""

### Example

"""sh
curl -X PUT "http://autocompletion:autocompletion@127.0.0.1:8080/index/twitter_accounts/store" -H "Content-Type: application/json" -d '{"id": "0001", "terms": ["@SamuelBerthe", "Samuel", "BERTHE"], "label": "Just put online a new side project: http://pipe.grep.so. Try it out !", "Payload": {"RT": 42, "favorites": 1337}, "score":42.0}'
"""

## Search

GET /index/#content-type/search?seed,score,pretty,payload

- **content-type**: "username", "article-title" or whatever...
- request **query**:
  - **seed** (string): Term to autocomplete (can include URL-encoded spaces)
  - **ordered** (bool, optional, default=true): return result ordered by score
  - **score** (bool, optional, default=false): for each result, return the score
  - **pretty** (bool, optional, default=false): for each label, add <strong></strong> tags around the terms that match with the seed parameter
  - **payload** (bool, optional, default=false): for each result, return the payload, if provided

### Return

"""js
{ results: [...] }
"""

### Example

"""sh
curl -X GET "http://127.0.0.1:8080/index/twitter_accounts/search?seed=sam%20ber&pretty=true&payload=true&score=true"
"""

## Search

GET /index/#content-type/#id?payload

- **content-type**: "username", "article-title" or whatever...
- **id**: content ID to retrieve
- request **query**:
  - **payload** (bool, optional, default=true): return the payload, if provided

### Return

"""js
{ doc: {...} }
"""

### Example

"""sh
curl -X GET "http://127.0.0.1:8080/index/twitter_accounts/00042
"""

# @TODO
- Add a tokenizer
- Add a english anayliser

