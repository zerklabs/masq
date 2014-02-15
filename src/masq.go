package main

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	mrand "math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

type HiddenMessage struct {
	Key      string `json:"key"`
	Url      string `json:"url"`
	Duration string `json:"duration"`
}

type ShowMessage struct {
	Value string `json:"value"`
}

type PasswordMessage struct {
	Password string `json:"password"`
}

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var pool *redis.Pool
var redisServer = flag.String("host", "127.0.0.1", "Redis Server")
var redisServerPort = flag.Int("port", 6379, "Redis Server Port")
var responseUrl = flag.String("url", "https://passwords.cobhamna.com", "Server Response URL")
var listenOn = flag.Int("listen", 8080, "Port to run the webserver on")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// bind the command line flags
	flag.Parse()

	redisUri := fmt.Sprintf("%s:%d", *redisServer, *redisServerPort)

	setupPool(redisUri)

	http.HandleFunc("/hide", hideHandler)
	http.HandleFunc("/show", showHandler)
	http.HandleFunc("/passwords", passwordsHandler)
	log.Printf("[Startup] Listening on: %d, Redis URI: %s", *listenOn, redisUri)
	http.ListenAndServe(fmt.Sprintf(":%d", *listenOn), nil)
}

func setupPool(server string) {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.Fatal(err)
			}

			return err
		},
	}
}

func hideHandler(w http.ResponseWriter, r *http.Request) {
	// set the contet type of the response to application/javascript
	w.Header().Set("Content-Type", "application/javascript")

	var duration string
	var data string
	expire := 24 * 3600
	inError := false

	// log the request
	logRequest(r)

	// generate a random key
	key := genRandomKey()

	// placeholder for storing data
	premadeUrl := url.Values{}

	premadeUrl.Set("key", key)

	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		handleError(w, "Failed to parse submit data", 400)

		inError = true
	}

	duration = r.Form.Get("duration")
	data = r.Form.Get("data")

	// if duration is not set, return a bad request
	if len(data) == 0 {
		handleError(w, "Missing Input (data)", 400)

		inError = true
	}

	switch duration {
	case "24":
		expire = 24 * 3600
		break
	case "48":
		expire = 48 * 3600
		break
	case "72":
		expire = 72 * 3600
		break
	case "1w":
		expire = 168 * 3600
	}

	conn := pool.Get()
	defer conn.Close()

	conn.Send("SET", key, data)
	conn.Send("EXPIRE", key, expire)
	conn.Flush()
	_, err := conn.Receive()

	if err != nil {
		log.Fatal(err)
		handleError(w, "Failed to retrieve value from Redis", 500)
		inError = true
	}

	if !inError {
		message := &HiddenMessage{
			Key:      key,
			Url:      fmt.Sprintf("%s/show?%s", *responseUrl, premadeUrl.Encode()),
			Duration: duration,
		}

		json.NewEncoder(w).Encode(message)
	}
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	// set the contet type of the response to application/javascript
	w.Header().Set("Content-Type", "application/javascript")

	var key string
	var data string
	inError := false

	// log the request
	logRequest(r)

	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
		handleError(w, "Failed to parse submit data", 400)

		inError = true
	}

	key = r.Form.Get("key")

	// if duration is not set, return a bad request
	if len(key) == 0 {
		handleError(w, "Missing Input (key)", 400)

		inError = true
	}

	conn := pool.Get()
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", key))

	if err != nil {
		log.Fatal(err)
		handleError(w, "Failed to retrieve value from Redis", 500)
		inError = true
	}

	if !inError {
		message := &ShowMessage{
			Value: data,
		}

		json.NewEncoder(w).Encode(message)
	}
}

/**
 *  masq-dev:dictionary is a zset
 */
func passwordsHandler(w http.ResponseWriter, r *http.Request) {
	// set the contet type of the response to application/javascript
	w.Header().Set("Content-Type", "application/javascript")

	inError := false

	mrand.Seed(time.Now().UTC().UnixNano())
	r1 := mrand.Intn(80000)
	r2 := mrand.Intn(80000)

	conn := pool.Get()
	defer conn.Close()

	// get first word of password
	w1, err := redis.Strings(conn.Do("ZRANGE", "masq-dev:dictionary", r1, r1))

	if err != nil {
		log.Fatal(err)
		handleError(w, "Failed to retrieve value from Redis", 500)
		inError = true
	}

	// get second word of password
	w2, err := redis.Strings(conn.Do("ZRANGE", "masq-dev:dictionary", r2, r2))

	if err != nil {
		log.Fatal(err)
		handleError(w, "Failed to retrieve value from Redis", 500)
		inError = true
	}

	randDigit := mrand.Intn(20000)

	// log the request
	logRequest(r)

	if !inError {
		message := &PasswordMessage{
			Password: fmt.Sprintf("%s%s%d", strings.Title(w1[0]), strings.Title(w2[0]), randDigit),
		}

		json.NewEncoder(w).Encode(message)
	}
}

func logRequest(r *http.Request) {
	log.Printf("[Request] %s %s - %s", r.Method, r.RequestURI, r.RemoteAddr)
}

func genRandomKey() string {
	c := 20
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	hash := sha1.New()

	hash.Write(b)

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func handleError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ErrorMessage{Code: code, Message: message})

	log.Printf("[Error in Request] %s", message)
}
