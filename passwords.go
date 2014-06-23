package main

import (
	"fmt"
	"log"
	mrand "math/rand"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	a "github.com/zerklabs/auburn-http"
)

// <prefix>:dictionary is a zset
func passwordsHandler(req a.HttpTransaction) {
	strong := req.Query("strong")

	if strong == "1" || strong == "true" {
		pwd := generateStrongPassword(32)

		println(pwd)

		req.RespondWithJSON(struct {
			Password string `json:"password"`
		}{
			Password: pwd,
		})

		return

	} else {
		conn := pool.Get()
		defer conn.Close()

		dictionaryKey := fmt.Sprintf("%s:dictionary", *redisKeyPrefix)
		mrand.Seed(time.Now().UTC().UnixNano())
		words, err := redis.Strings(conn.Do("SRANDMEMBER", dictionaryKey, 2))

		if err != nil {
			log.Println(err)
			req.Error("Failed to retrieve value from Redis", 500)
			return
		}

		if len(words) == 0 {
			req.Error("Failed to find value in Redis list", 404)
			return
		}

		randDigit := mrand.Intn(20000)

		req.RespondWithJSON(struct {
			Password string `json:"password"`
		}{
			Password: fmt.Sprintf("%s%s%d", strings.Title(words[0]), strings.Title(words[1]), randDigit),
		})
	}
}
