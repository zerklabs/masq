package main

import (
	"flag"

	"github.com/garyburd/redigo/redis"
	"github.com/zerklabs/auburn/http"
	"github.com/zerklabs/auburn/redis"
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
	)

	// bind the command line flags
	flag.Parse()

	pool = auburnredis.NewPool(*redisAddress, "")
	server, err := http.NewServer()

	if err != nil {
		panic(err)
	}

	server.Options.EnableLogging = true
	server.Options.EnableCors = true

	server.AddRouteForMethod("/2/hide", http.POST, hideHandler)
	server.AddRouteForMethod("/hide", http.POST, hideHandler)

	server.AddRouteForMethod("/2/show", http.GET, showHandler)
	server.AddRouteForMethod("/show", http.GET, showHandler)

	server.AddRouteForMethod("/2/passwords", http.GET, passwordsHandler)
	server.AddRouteForMethod("/passwords", http.GET, passwordsHandler)

	server.Start()
}
