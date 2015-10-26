package server

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/pat"
)

var (
	router *pat.Router
)

func init() {
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

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	ss := http.StripPrefix("", http.FileServer(http.Dir("./templates/")))
	r.PathPrefix("/").Handler(ss)

	return r
}

