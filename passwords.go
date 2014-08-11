package main

import (
	"fmt"
	mrand "math/rand"
	"strings"
	"time"

	log "code.google.com/p/log4go"
	"github.com/garyburd/redigo/redis"
	a "github.com/zerklabs/auburn-http"
)

// <prefix>:dictionary is a zset
func passwordsHandler(req *a.HttpTransaction) (error, int) {
	strong := req.Query("strong")

	if strong == "1" || strong == "true" {
		pwd := generateStrongPassword(32)

		return req.RespondWithJSON(struct {
			Password string `json:"password"`
		}{
			Password: pwd,
		})

	} else {
		conn := pool.Get()
		defer conn.Close()

		dictionaryKey := fmt.Sprintf("%s:dictionary", *redisKeyPrefix)
		mrand.Seed(time.Now().UTC().UnixNano())
		words, err := redis.Strings(conn.Do("SRANDMEMBER", dictionaryKey, 2))

		if err != nil {
			log.Error(err)
			return req.Error("Failed to retrieve value from Redis", 500)
		}

		if len(words) == 0 {
			return req.Error("Failed to find value in Redis list", 404)
		}

		randDigit := mrand.Intn(20000)

		return req.RespondWithJSON(struct {
			Password string `json:"password"`
		}{
			Password: fmt.Sprintf("%s%s%d", strings.Title(words[0]), strings.Title(words[1]), randDigit),
		})
	}
}