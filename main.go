package main

import (
	"os"
	"imgresizer/server"
)

//Host eg: localhost:3001 ; or website.com
var Port = os.Getenv("PORT")

func main() {
	if Port == "" {
		Port = "3666"
	}
	server.StartServer(Port)
}
