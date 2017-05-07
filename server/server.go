package server

import (
	"fmt"
	"github.com/gorilla/pat"
	cache "github.com/pmylund/go-cache"
	"log"
	"net/http"

	"errors"
	"gopkg.in/redis.v3"
	"os"
	"time"
)

var (
	router *pat.Router
	Cache  *cache.Cache
	Client *redis.Client

	//	Client *redis.ClusterClient
)

var REDIS_HOST = os.Getenv("REDIS_HOST")
var REDIS_PASS = os.Getenv("REDIS_PASSWORD")
var REDIS_PORT = os.Getenv("REDIS_PORT")

var RedisCacheDisabled = os.Getenv("rediscachedisabled")

func init() {

	if RedisCacheDisabled == "true" {
		log.Println("INFO: Redis is disabled, wont connect !!")
	} else {

		if REDIS_HOST == "" {
			REDIS_HOST = "g7-box"
			log.Println("WARN: env REDIS_HOST is not set, will use default ", REDIS_HOST)
		}

		if REDIS_PASS == "" {
			log.Println("WARN: env REDIS_PASS is not set, was this intentional?")
		}

		if REDIS_PORT == "" {
			REDIS_PORT = "6379"
			log.Println("WARN: env REDIS_PORT is not set, will use default ", REDIS_PORT)
		}

		Cache = cache.New(30 * time.Minute, 60 * time.Second)

		Client = redis.NewClient(&redis.Options{
			Addr:     REDIS_HOST + ":" + REDIS_PORT,
			Password: REDIS_PASS, // password set
			DB:       0, // use default DB
		})
		//	addrs :=  strings.Split(REDIS_HOST,",")
		//	log.Println("karai")
		//	log.Println(addrs)
		//	Client = redis.NewClusterClient(&redis.ClusterOptions{
		////		Addrs: []string{"82.196.8.72:7000", "146.185.154.216:7000", "82.196.9.79:7000"},
		//		Addrs: addrs,
		//		Password: REDIS_PASS, // no password set
		//	})
		pong, err := Client.Ping().Result()
		log.Println("Redis initial Ping result -> ", pong, err)

	}

}

//NewServer return pointer to new created server object
func NewServer(Port string) *http.Server {
	router = InitRouting()
	return &http.Server{
		Addr:    ":" + Port,
		Handler: router,
	}
}

//StartServer start and listen @server
func StartServer(Port string) {
	log.Println("Starting server")
	s := NewServer(Port)
	fmt.Println("Server starting --> " + Port)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalln("Error: %v", err)
	}
}

func InitRouting() *pat.Router {

	r := pat.New()

	r.Get("/img/", ImgHandler)
	r.Post("/img/", ImgPostHandler)
	r.Post("/r/img/", RedisImgPostHandler)
	r.Post("/meme/", MemeHandler)
	//r.Post("/gif/", GifPostHandler)

	//	r.Get("/imgbykey/{operator}/{key}/{value}/", RedisImgGetHandler)

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	ss := http.StripPrefix("", http.FileServer(http.Dir("./templates/")))
	r.PathPrefix("/").Handler(ss)

	return r
}

func GetFromCache(key string) (string, bool) {
	var valid bool = false
	if RedisCacheDisabled == "true" {
		log.Println("Redis Cache is Disabled")
		return "", valid
	}
	cached, err := Client.Get(key).Result()
	//
	if err == redis.Nil {
		log.Println("Error:  GetFromCache -> ", key, " does not exists")
	} else if err != nil {
		log.Println("Error: ", err)
	} else {
		//log.Println("Serve from cache")
		valid = true
	}
	return cached, valid
}

func SetToCache(key string, value interface{}, expiration time.Duration) error {
	var valid error = nil
	if RedisCacheDisabled == "true" {
		log.Println("Redis Cache is Disabled")
		valid = errors.New("Redis Cache is Disabled")
	} else {
		Client.Set(key, value, expiration)
	}
	return valid
}
