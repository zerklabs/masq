package main

import (
	"flag"
	"fmt"
	"github.com/cabrel/auburn"
	"github.com/garyburd/redigo/redis"
	"log"
	mrand "math/rand"
	"net/url"
	"runtime"
	"strings"
	"time"
)

var pool *redis.Pool
var redisServer = flag.String("host", "127.0.0.1", "Redis Server")
var redisServerPort = flag.Int("port", 6379, "Redis Server Port")
var responseUrl = flag.String("url", "https://passwords.cobhamna.com", "Server Response URL")
var listenOn = flag.Int("listen", 8080, "Port to run the webserver on")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// bind the command line flags
	flag.Parse()

	redisUri := fmt.Sprintf("%s:%d", *redisServer, *redisServerPort)

	setupPool(redisUri)

	server := &auburn.AuburnHttpServer{HttpPort: *listenOn}

	server.Handle("/hide", hideHandler)
	server.Handle("/show", showHandler)
	server.Handle("/passwords", passwordsHandler)
	server.Start()
}

//
func setupPool(server string) {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.Fatal(err)
			}

			return err
		},
	}
}

//
func hideHandler(req *auburn.AuburnHttpRequest) {
	expire := 24 * 3600

	// generate a random key
	key := auburn.genRandomKey()

	// placeholder for storing data
	premadeUrl := url.Values{}
	premadeUrl.Set("key", key)

	duration, err := req.GetValue("duration")

	if err != nil {
		req.Error("Failed to get `duration` from Form", 400)
	}

	data, err := req.GetValue("data")

	if err != nil {
		req.Error("Failed to get `data` from Form", 400)
	}

	if len(data) == 0 {
		req.Error("Missing `data` value", 400)
	}

	switch duration {
	case "24":
		expire = 24 * 3600
		break
	case "48":
		expire = 48 * 3600
		break
	case "72":
		expire = 72 * 3600
		break
	case "1w":
		expire = 168 * 3600
	}

	conn := pool.Get()
	defer conn.Close()

	conn.Send("SET", key, data)
	conn.Send("EXPIRE", key, expire)
	conn.Flush()

	req.Respond(struct {
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
	key, err := req.GetValue("key")

	if err != nil {
		req.Error("Failed to get `key` from Form", 400)
	}

	conn := pool.Get()
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", key))

	if err != nil {
		req.Error("Failed to retrieve value from Redis", 500)
	}

	req.Respond(struct {
		Value string `json:"value"`
	}{
		Value: data,
	})
}

// masq-dev:dictionary is a zset
func passwordsHandler(req *auburn.AuburnHttpRequest) {
	mrand.Seed(time.Now().UTC().UnixNano())
	r1 := mrand.Intn(80000)
	r2 := mrand.Intn(80000)

	conn := pool.Get()
	defer conn.Close()

	// get first word of password
	w1, err := redis.Strings(conn.Do("ZRANGE", "masq-dev:dictionary", r1, r1))

	if err != nil {
		req.Error("Failed to retrieve value from Redis", 500)
	}

	// get second word of password
	w2, err := redis.Strings(conn.Do("ZRANGE", "masq-dev:dictionary", r2, r2))

	if err != nil {
		req.Error("Failed to retrieve value from Redis", 500)
	}

	randDigit := mrand.Intn(20000)

	req.Respond(struct {
		Password string `json:"password"`
	}{
		Password: fmt.Sprintf("%s%s%d", strings.Title(w1[0]), strings.Title(w2[0]), randDigit),
	})
}
