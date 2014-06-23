package main

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	a "github.com/zerklabs/auburn-http"
)

//
func showHandler(req a.HttpTransaction) {
	conn := pool.Get()
	defer conn.Close()

	key := req.Query("key")

	if len(key) == 0 {
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
