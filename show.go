package main

import (
	"errors"
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
		return response.Error(400, errors.New("Failed to get key from Form"))
	}

	if key == "dictionary" {
		return response.Error(401, errors.New("Invalid request type"))
	}

	uniqueKey := fmt.Sprintf("%s:%s", *redisKeyPrefix, key)
	data, err := redis.String(cluster.Do("GET", uniqueKey))
	if err != nil {
		return response.JsonError(400, "Failed to retrieve value")
	}

	return req.Json(struct {
		Value string `json:"value"`
	}{
		Value: data,
	})
}
