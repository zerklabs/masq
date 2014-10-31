package main

import (
	"fmt"
	"net/url"

	"code.google.com/p/go.net/context"

	"github.com/zerklabs/auburn/http"
	"github.com/zerklabs/auburn/http/response"
	"github.com/zerklabs/auburn/utils"
)

//
func hideHandler(ctx context.Context, req http.HttpTransaction) response.Response {
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
		return req.Error(400, "Missing `data` value")
	}

	uniqueKey := fmt.Sprintf("%s:%s", *redisKeyPrefix, key)

	_, err := cluster.Do("SET", uniqueKey, data)
	if err != nil {
		http.Log.Error(err)
		return req.Error(500, err.Error())
	}

	// include consideration for no duration
	if durations[duration] > 0 {
		_, err := cluster.Do("EXPIRE", uniqueKey, durations[duration])
		if err != nil {
			http.Log.Error(err)
			return req.Error(500, err.Error())
		}
	}

	http.Log.Infof("%s expires in %v", uniqueKey, durations[duration])

	return req.Json(struct {
		Key      string `json:"key"`
		Url      string `json:"url"`
		Duration string `json:"duration"`
	}{
		Key:      key,
		Url:      fmt.Sprintf("%s/show?%s", *responseUrl, premadeUrl.Encode()),
		Duration: duration,
	})
}
