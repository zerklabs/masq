package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	a "github.com/zerklabs/auburn-http"
)

//
func showHandler(req *a.HttpTransaction) (error, int) {
	conn := pool.Get()
	defer conn.Close()

	key := req.Query("key")

	if len(key) == 0 {
		return req.Error("Failed to get `key` from Form", 400)
	}

	if key == "dictionary" {
		return req.Error("Invalid Request", 400)
	}

	uniqueKey := fmt.Sprintf("%s:%s", *redisKeyPrefix, key)
	data, err := redis.String(conn.Do("GET", uniqueKey))

	if err != nil {
		a.Log.Error(err)
		return req.Error("Failed to retrieve value from Redis", 500)
	}

	return req.RespondWithJSON(struct {
		Value string `json:"value"`
	}{
		Value: data,
	})
}
