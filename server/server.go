package server

import (
	"code_runner/workers"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Server struct {
	Server *http.ServeMux
}

func InitServer() (server Server) {
	server = Server{http.NewServeMux()}
	return server
}

func (server Server) Run() {
	server.addRouteHandlers()
}

func (server Server) addRouteHandlers() {
	server.Server.HandleFunc("/api/run", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		job := workers.Job{Language: workers.Python, Code: "for i in range(5):\n\tprint('yooo')"}
		codeRunner, err := workers.CreateCodeRunner(job)

		if err != nil {
			slog.Error(err.Error())
		}

		outputChannel := make(chan workers.JobResult)
		go codeRunner.RunCode(outputChannel)
		jobResult := <-outputChannel

		err = json.NewEncoder(writer).Encode(jobResult)
		if err != nil {
			slog.Error(err.Error())
		}
	})
}
