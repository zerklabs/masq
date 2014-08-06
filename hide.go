package main

import (
	"fmt"
	"net/url"

	a "github.com/zerklabs/auburn-http"
)

//
func hideHandler(req a.HttpTransaction) {
	conn := pool.Get()
	defer conn.Close()

	// generate a random key
	key := a.GenRandomKey()

	// placeholder for storing data
	premadeUrl := url.Values{}
	premadeUrl.Set("key", key)

	duration := req.Query("duration")

	if len(duration) == 0 {
		duration = "24h"
	}

	data := req.Query("data")

	if len(data) == 0 {
		req.Error("Missing `data` value", 400)
		return
	}

	uniqueKey := fmt.Sprintf("%s:%s", *redisKeyPrefix, key)

	conn.Send("SET", uniqueKey, data)

	// include consideration for no duration
	if durations[duration] > 0 {
		conn.Send("EXPIRE", uniqueKey, durations[duration])
	}

	conn.Flush()

	req.RespondWithJSON(struct {
		Key      string `json:"key"`
		Url      string `json:"url"`
		Duration string `json:"duration"`
	}{
		Key:      key,
		Url:      fmt.Sprintf("%s/show?%s", *responseUrl, premadeUrl.Encode()),
		Duration: duration,
	})
}
