package main

import (
	"errors"
	"fmt"
	mrand "math/rand"
	"strings"
	"time"

	"code.google.com/p/go.net/context"

	"github.com/garyburd/redigo/redis"
	"github.com/zerklabs/auburn/http"
	"github.com/zerklabs/auburn/http/response"
	"github.com/zerklabs/auburn/log"
)

// <prefix>:dictionary is a zset
func passwordsHandler(ctx context.Context, req http.HttpTransaction) response.Response {
	strong := req.Query("strong")

	if strong == "1" || strong == "true" {
		pwd := generateStrongPassword(32)

		return req.Json(struct {
			Password string `json:"password"`
		}{
			Password: pwd,
		})

	} else {
		dictionaryKey := fmt.Sprintf("%s:dictionary", *redisKeyPrefix)
		mrand.Seed(time.Now().UTC().UnixNano())
		words, err := redis.Strings(cluster.Do("SRANDMEMBER", dictionaryKey, 2))
		if err != nil {
			log.Error(err)
			return response.Error(500, errors.New("Failed to retrieve value"))
		}
		if len(words) == 0 {
			return response.Error(500, errors.New("Failed to find value in list"))
		}

		randDigit := mrand.Intn(20000)

		return req.Json(struct {
			Password string `json:"password"`
		}{
			Password: fmt.Sprintf("%s%s%d", strings.Title(words[0]), strings.Title(words[1]), randDigit),
		})
	}
}
