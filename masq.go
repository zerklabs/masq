package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/zerklabs/auburn"
	"log"
	mrand "math/rand"
	"net/url"
	"runtime"
	"strings"
	"time"
)

var (
	redisServer     = flag.String("redisip", "127.0.0.1", "Redis Server")
	redisServerPort = flag.Int("redisport", 6379, "Redis Server Port")
	responseUrl     = flag.String("url", "https://passwords.cobhamna.com", "Server Response URL")
	redisKeyPrefix  = flag.String("prefix", "masq", "Key prefix in Redis")
	listenIP        = flag.String("host", "127.0.0.1", "Port to run the webserver on")
	listenOn        = flag.Int("listen", 8080, "Port to run the webserver on")

	redisUri string

	// predefined string -> int (as seconds) durations
	durations = map[string]int{
		"5m":  5 * 60,
		"10m": 10 * 60,
		"15m": 15 * 60,
		"30m": 30 * 60,
		"1h":  3600,
		"24h": 24 * 3600,
		"48h": 48 * 3600,
		"72h": 72 * 3600,
		"1w":  168 * 3600,
	}
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// bind the command line flags
	flag.Parse()

	redisUri = fmt.Sprintf("%s:%d", *redisServer, *redisServerPort)

	server := auburn.New(*listenIP, *listenOn, "", "")

	server.Handle("/2/hide", hideHandler)
	server.Handle("/2/show", showHandler)
	server.Handle("/2/passwords", passwordsHandler)
	server.Start()
}

//
func hideHandler(req *auburn.AuburnHttpRequest) {
	conn, err := redis.Dial("tcp", redisUri)

	if err != nil {
		req.Error("Failed to connect to redis", 500)
	}

	defer conn.Close()

	// generate a random key
	key := auburn.GenRandomKey()

	// placeholder for storing data
	premadeUrl := url.Values{}
	premadeUrl.Set("key", key)

	duration, err := req.GetValue("duration")

	if err != nil {
		log.Print(err)
		req.Error("Failed to get `duration` from Form", 400)
	}

	if len(duration) == 0 {
		duration = "24h"
	}

	data, err := req.GetValue("data")

	if err != nil {
		log.Print(err)
		req.Error("Failed to get `data` from Form", 400)
	}

	if len(data) == 0 {
		req.Error("Missing `data` value", 400)
	}

	uniqueKey := fmt.Sprintf("%s:%s", *redisKeyPrefix, key)

	conn.Send("SET", uniqueKey, data)
	conn.Send("EXPIRE", uniqueKey, durations[duration])
	conn.Flush()

	req.RespondWithJSON(struct {
		Key      string `json:"key"`
		Url      string `json:"url"`
		Duration string `json:"duration"`
	}{
		Key:      key,
		Url:      fmt.Sprintf("%s/show?%s", *responseUrl, premadeUrl.Encode()),
		Duration: duration,
	})
}

//
func showHandler(req *auburn.AuburnHttpRequest) {
	conn, err := redis.Dial("tcp", redisUri)

	if err != nil {
		req.Error("Failed to connect to redis", 500)
	}

	defer conn.Close()

	key, err := req.GetValue("key")

	if err != nil {
		req.Error("Failed to get `key` from Form", 400)
	}

	if key == "dictionary" {
		req.Error("Invalid Request", 401)
	}

	uniqueKey := fmt.Sprintf("%s:%s", *redisKeyPrefix, key)
	data, err := redis.String(conn.Do("GET", uniqueKey))

	if err != nil {
		log.Print(err)
		req.Error("Failed to retrieve value from Redis", 500)
	}

	req.RespondWithJSON(struct {
		Value string `json:"value"`
	}{
		Value: data,
	})
}

// masq-dev:dictionary is a zset
func passwordsHandler(req *auburn.AuburnHttpRequest) {
	conn, err := redis.Dial("tcp", redisUri)

	if err != nil {
		req.Error("Failed to connect to redis", 500)
	}

	defer conn.Close()

	dictionaryKey := fmt.Sprintf("%s:dictionary", *redisKeyPrefix)

	mrand.Seed(time.Now().UTC().UnixNano())
	r1 := mrand.Intn(80000)
	r2 := mrand.Intn(80000)

	// get first word of password
	w1, err := redis.Strings(conn.Do("ZRANGE", dictionaryKey, r1, r1))

	if err != nil {
		req.Error("Failed to retrieve value from Redis", 500)
	}

	// get second word of password
	w2, err := redis.Strings(conn.Do("ZRANGE", dictionaryKey, r2, r2))

	if err != nil {
		req.Error("Failed to retrieve value from Redis", 500)
	}

	if len(w1) == 0 || len(w2) == 0 {
		req.Error("Failed to find value in Redis list", 404)
	}

	randDigit := mrand.Intn(20000)

	req.RespondWithJSON(struct {
		Password string `json:"password"`
	}{
		Password: fmt.Sprintf("%s%s%d", strings.Title(w1[0]), strings.Title(w2[0]), randDigit),
	})
}
