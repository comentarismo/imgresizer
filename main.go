package main

import (
	"imgresizer/server"
	"os"
)

//Host eg: localhost:3001 ; or website.com
var Port = os.Getenv("PORT")

func main() {
	if Port == "" {
		Port = "3666"
	}
	server.StartServer(Port)
}
