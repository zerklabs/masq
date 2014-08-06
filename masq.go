package main

import (
	"flag"
	"time"

	"github.com/garyburd/redigo/redis"
	a "github.com/zerklabs/auburn-http"
)

var (
	// predefined string -> int (as seconds) durations
	durations = map[string]int{
		"none": 0,
		"5m":   5 * 60,
		"10m":  10 * 60,
		"15m":  15 * 60,
		"30m":  30 * 60,
		"1h":   3600,
		"24h":  24 * 3600,
		"48h":  48 * 3600,
		"72h":  72 * 3600,
		"1w":   168 * 3600,
	}

	pool *redis.Pool

	responseUrl    = flag.String("url", "https://passwords.cobhamna.com", "Server Response URL")
	redisKeyPrefix = flag.String("prefix", "masq", "Key prefix in Redis")
)

func main() {
	var (
		redisAddress = flag.String("redis-address", "127.0.0.1:6379", "Redis Server")
		httpAddress  = flag.String("http-address", "127.0.0.1:8080", "IP:PORT to run the webserver on")
	)

	// bind the command line flags
	flag.Parse()

	pool = newPool(*redisAddress)
	server := a.New(*httpAddress)

	server.AddRouteForMethod("/2/hide", a.POST, hideHandler)
	server.AddRouteForMethod("/hide", a.POST, hideHandler)

	server.AddRouteForMethod("/2/show", a.GET, showHandler)
	server.AddRouteForMethod("/show", a.GET, showHandler)

	server.AddRouteForMethod("/2/passwords", a.GET, passwordsHandler)
	server.AddRouteForMethod("/passwords", a.GET, passwordsHandler)

	server.Start()
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)

			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
