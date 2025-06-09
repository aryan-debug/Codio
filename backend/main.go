package main

import (
	"code_runner/server"
	"log"
	"net/http"
)

func main() {
	server := server.InitServer()
	server.Run()
	err := http.ListenAndServe(":8080", server.Server)
	if err != nil {
		log.Fatalf("Error starting the server: %s", err)
	}
}
