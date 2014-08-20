package main

import (
	"fmt"
	"net/url"

	"github.com/zerklabs/auburn/http"
	"github.com/zerklabs/auburn/utils"
)

//
func hideHandler(req *http.HttpTransaction) (error, int) {
	conn := pool.Get()
	defer conn.Close()

	// generate a random key
	key, _ := utils.GenRandomKey()

	// placeholder for storing data
	premadeUrl := url.Values{}
	premadeUrl.Set("key", key)

	duration := req.Query("duration")

	if len(duration) == 0 {
		duration = "24h"
	}

	data := req.Query("data")

	if len(data) == 0 {
		return req.Error("Missing `data` value", 400)
	}

	uniqueKey := fmt.Sprintf("%s:%s", *redisKeyPrefix, key)

	conn.Send("SET", uniqueKey, data)

	// include consideration for no duration
	if durations[duration] > 0 {
		conn.Send("EXPIRE", uniqueKey, durations[duration])
	}

	http.Log.Infof("%s expires in %v", uniqueKey, durations[duration])

	conn.Flush()

	return req.RespondWithJSON(struct {
		Key      string `json:"key"`
		Url      string `json:"url"`
		Duration string `json:"duration"`
	}{
		Key:      key,
		Url:      fmt.Sprintf("%s/show?%s", *responseUrl, premadeUrl.Encode()),
		Duration: duration,
	})
}
