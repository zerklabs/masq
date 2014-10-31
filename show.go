package main

import (
	"fmt"

	"code.google.com/p/go.net/context"

	"github.com/garyburd/redigo/redis"
	"github.com/zerklabs/auburn/http"
	"github.com/zerklabs/auburn/http/response"
)

//
func showHandler(ctx context.Context, req http.HttpTransaction) response.Response {
	key := req.Query("key")

	if len(key) == 0 {
		return req.Error(400, "Failed to get `key` from Form")
	}

	if key == "dictionary" {
		return req.Error(400, "Invalid Request")
	}

	uniqueKey := fmt.Sprintf("%s:%s", *redisKeyPrefix, key)
	data, err := redis.String(cluster.Do("GET", uniqueKey))
	if err != nil {
		http.Log.Error(err)
		return req.Error(500, "Failed to retrieve value from Redis")
	}

	return req.Json(struct {
		Value string `json:"value"`
	}{
		Value: data,
	})
}
