package server

import (
	"code_runner/workers"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
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
	server.Server.HandleFunc("/api/run", runCode)
}

// Reads the language and code from the request and spawns a new job to execute the code
func runCode(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var job workers.Job
	err := json.NewDecoder(req.Body).Decode(&job)
	if err != nil {
		panic(err)
	}
	jobResult := runJob(workers.Job{Language: job.Language, Code: job.Code})
	content, _ := (json.Marshal(jobResult))
	os.Stdout.Write(content)
	err = json.NewEncoder(writer).Encode(jobResult)
	if err != nil {
		slog.Error(err.Error())
	}
}

// Creates a code runner, gives it a job and waits for it to finish
func runJob(job workers.Job) workers.JobResult {
	codeRunner, err := workers.CreateCodeRunner(job)

	if err != nil {
		slog.Error(err.Error())
		return workers.JobResult{Output: "", Error: err.Error()}
	}

	outputChannel := make(chan workers.JobResult)
	go codeRunner.RunCode(outputChannel)
	jobResult := <-outputChannel
	return jobResult
}
