package main

import (
	"code_runner/server"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	server := server.InitServer()
	server.Run()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
	})

	err := http.ListenAndServe(":8080", c.Handler(server.Server))
	if err != nil {
		log.Fatalf("Error starting the server: %s", err)
	}
}
