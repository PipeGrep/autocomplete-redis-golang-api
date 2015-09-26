package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
)

type Config_Redis struct {
	Hostname string
	Port     interface{}
	Password string
}

type Config_Basic struct {
	Username string
	Password string
}

type Config_HttpServer struct {
	Port  interface{}
	Basic Config_Basic
}

type Config_Index struct {
	MinSizeToken int
}

type Config struct {
	Redis Config_Redis
	Http  Config_HttpServer
	Index Config_Index
}

var config = Config{}

func init() {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		log.Fatal("Cannot open config.json: " + e.Error())
	}

	e = json.Unmarshal(file, &config)
	if e != nil {
		log.Fatal("Parsing error, config.json is not a valid Json file: " + e.Error())
	}

	init_from_env(&config)
}

// get env variables
func init_from_env(config *Config) {
	// REDIS.HOSTNAME
	if config.Redis.Hostname[0] == '$' {
		config.Redis.Hostname = os.Getenv(config.Redis.Hostname[1:])
	}

	// REDIS.PORT
	if reflect.TypeOf(config.Redis.Port).String() == "string" {
		config.Redis.Port, _ = strconv.Atoi(os.Getenv(config.Redis.Port.(string)[1:]))
	} else if reflect.TypeOf(config.Redis.Port).String() == "float64" {
		config.Redis.Port = int(config.Redis.Port.(float64))
	}
	if reflect.TypeOf(config.Redis.Port).String() != "int" {
		log.Fatal("Incorrect type in configuration file: redis.port.")
	}

	// REDIS.PASSWORD
	if config.Redis.Password != "" && config.Redis.Password[0] == '$' {
		config.Redis.Password = os.Getenv(config.Redis.Password[1:])
	}

	// HTTP.PORT
	if reflect.TypeOf(config.Http.Port).String() == "string" {
		config.Http.Port, _ = strconv.Atoi(os.Getenv(config.Http.Port.(string)[1:]))
	} else if reflect.TypeOf(config.Http.Port).String() == "float64" {
		config.Http.Port = int(config.Http.Port.(float64))
	}
	if reflect.TypeOf(config.Http.Port).String() != "int" {
		log.Fatal("Incorrect type in configuration file: http.port.")
	}

	// HTTP.USERNAME
	if config.Http.Basic.Username != "" && config.Http.Basic.Username[0] == '$' {
		config.Http.Basic.Username = os.Getenv(config.Http.Basic.Username[1:])
	}

	// HTTP.PASSWORD
	if config.Http.Basic.Password != "" && config.Http.Basic.Password[0] == '$' {
		config.Http.Basic.Password = os.Getenv(config.Http.Basic.Password[1:])
	}
}
