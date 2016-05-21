package server

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/pat"
	cache "github.com/pmylund/go-cache"

	"time"
	"gopkg.in/redis.v3"
	"os"
)

var (
	router *pat.Router
	Cache *cache.Cache
	Client *redis.Client
//	Client *redis.ClusterClient
)

var REDIS_HOST = os.Getenv("REDIS_HOST")
var REDIS_PASS = os.Getenv("REDIS_PASSWORD")
var REDIS_PORT = os.Getenv("REDIS_PORT")

func init() {
	if REDIS_HOST == "" {
		log.Fatal("env REDIS_HOST is not set")
	}

	if REDIS_PASS == "" {
		log.Println("env REDIS_PASS is not set, was this intentional?")
	}

	if REDIS_PORT == "" {
		log.Fatal("env REDIS_PORT is not set")
	}

	Cache = cache.New(30*time.Minute, 60*time.Second)

	Client = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST+":"+REDIS_PORT,
		Password: REDIS_PASS, // password set
		DB:       0,  // use default DB
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
	log.Println(pong, err)
}

//NewServer return pointer to new created server object
func NewServer(Port string) *http.Server {
	router = InitRouting()
	return &http.Server{
		Addr:    ":"+Port,
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
//	r.Get("/imgbykey/{operator}/{key}/{value}/", RedisImgGetHandler)

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	ss := http.StripPrefix("", http.FileServer(http.Dir("./templates/")))
	r.PathPrefix("/").Handler(ss)

	return r
}


