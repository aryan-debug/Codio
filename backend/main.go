package main

import (
	"code_runner/server"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	server := server.InitServer()
	server.Run()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("CLIENT_URL")},
	})

	err = http.ListenAndServe(":8080", c.Handler(server.Server))
	if err != nil {
		log.Fatalf("Error starting the server: %s", err)
	}
}
