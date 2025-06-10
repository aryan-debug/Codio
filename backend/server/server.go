package server

import (
	"code_runner/workers"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type jobRequest struct {
	Language string
	Code     string
}

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

func runCode(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var job jobRequest
	err := json.NewDecoder(req.Body).Decode(&job)
	if err != nil {
		panic(err)
	}
	fmt.Println(job.Code)
	fmt.Println(job.Language)
	jobResult := runJob(workers.Job{Language: workers.LanguageMap[job.Language], Code: job.Code})
	fmt.Println(jobResult.Output)
	fmt.Println(jobResult.Error)
	content, _ := (json.Marshal(jobResult))
	os.Stdout.Write(content)
	err = json.NewEncoder(writer).Encode(jobResult)
	if err != nil {
		slog.Error(err.Error())
	}
}

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
