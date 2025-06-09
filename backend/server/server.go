package server

import (
	"code_runner/workers"
	"encoding/json"
	"fmt"
	"io"
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

func runCode(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	language := req.FormValue("language")
	file_content := getFileContent(req, "file")
	jobResult := runJob(workers.Job{Language: workers.LanguageMap[language], Code: string(file_content)})
	fmt.Println(jobResult.Output)
	fmt.Println(jobResult.Error)
	content, _ := (json.Marshal(jobResult))
	os.Stdout.Write(content)
	err := json.NewEncoder(writer).Encode(jobResult)
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

func getFileContent(req *http.Request, key string) string {
	file, _, err := req.FormFile(key)

	if err != nil {
		slog.Error(err.Error())
	}

	file_content, err := io.ReadAll(file)

	if err != nil {
		slog.Error(err.Error())
	}

	return string(file_content)
}
