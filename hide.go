package main

import (
	"errors"
	"fmt"
	"net/url"

	"code.google.com/p/go.net/context"

	"github.com/zerklabs/auburn/http"
	"github.com/zerklabs/auburn/http/response"
	"github.com/zerklabs/auburn/log"
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
		log.Error("missing data value")
		return response.Error(400, errors.New("Missing data value"))
	}

	uniqueKey := fmt.Sprintf("%s:%s", *redisKeyPrefix, key)

	_, err := cluster.Do("SET", uniqueKey, data)
	if err != nil {
		log.Error(err)
		return response.Error(500, err)
	}

	// include consideration for no duration
	if durations[duration] > 0 {
		_, err := cluster.Do("EXPIRE", uniqueKey, durations[duration])
		if err != nil {
			log.Error(err)
			return response.Error(500, err)
		}
	}

	log.Infof("%s expires in %v", uniqueKey, durations[duration])

	return req.Json(struct {
		Key      string `json:"key"`
		Url      string `json:"url"`
		Duration string `json:"duration"`
	}{
		Key:      key,
		Url:      fmt.Sprintf("%s/#/show?%s", *responseUrl, premadeUrl.Encode()),
		Duration: duration,
	})
}
