package server

import "net/http"

type Server struct {
	server *http.ServeMux
}

func InitServer() Server {
	return Server{http.NewServeMux()}
}
