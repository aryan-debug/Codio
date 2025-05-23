package workers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
)

var ResultsChannel = make(chan JobResult, 100)

type codeRunner struct {
	ID               int
	Work             chan Job
	WorkerQueue      chan chan Job
	QuitChan         chan bool
	containerManager *containerManager
}

func CreateCodeRunner(id int, workerQueue chan chan Job) (*codeRunner, error) {
	containerManager, err := CreateContainerManager()

	if err != nil {
		return nil, err
	}

	return &codeRunner{
		ID:               id,
		Work:             make(chan Job),
		WorkerQueue:      workerQueue,
		QuitChan:         make(chan bool),
		containerManager: containerManager,
	}, nil
}

func (cr *codeRunner) Start() {
	go func() {
		for {
			cr.WorkerQueue <- cr.Work

			select {
			case job := <-cr.Work:
				cr.RunCode(job)
			case <-cr.QuitChan:
				return
			}
		}
	}()
}

func (cr *codeRunner) Stop() {
	go func() {
		cr.QuitChan <- true
	}()
}

func (cr *codeRunner) RunCode(job Job) {
	filename := "test.py"
	if job.Language == Java {
		filename = "Main.java"
	}

	tempDirPath, err := writeToTempFile(job.Code, filename)

	if err != nil {
		ResultsChannel <- JobResult{ID: job.ID, Output: "", Error: err}
		return
	}

	defer os.RemoveAll(tempDirPath)

	imageName := LanguageImageMap[job.Language]

	resp, err := cr.containerManager.CreateContainer(imageName, tempDirPath)

	if err != nil {
		ResultsChannel <- JobResult{ID: job.ID, Output: "", Error: err}
		return
	}

	err = cr.containerManager.StartContainer(resp.ID)

	if err != nil {
		ResultsChannel <- JobResult{ID: job.ID, Output: "", Error: err}
		return
	}

	cr.containerManager.WaitForContainer(resp.ID, func(sc <-chan container.WaitResponse, ec <-chan error) {
		<-sc
	})

	stdout, stderr, err := cr.containerManager.GetContainerOutputParsed(resp.ID)

	if err != nil {
		ResultsChannel <- JobResult{ID: job.ID, Output: "", Error: err}
		return
	}

	err = cr.containerManager.RemoveContainer(resp.ID)

	if err != nil {
		ResultsChannel <- JobResult{ID: job.ID, Output: "", Error: err}
		return
	}

	var output string
	if stdout != "" && stderr != "" {
		output = fmt.Sprintf("STDOUT:\n%s\nSTDERR:\n%s", stdout, stderr)
	} else if stdout != "" {
		output = stdout
	} else if stderr != "" {
		output = stderr
	}

	ResultsChannel <- JobResult{ID: job.ID, Output: output, Error: nil}
}

func writeToTempFile(code string, filename string) (string, error) {
	tempDir, err := os.MkdirTemp("", "temp")

	if err != nil {
		return "", err
	}

	err = os.WriteFile(filepath.Join(tempDir, filename), []byte(code), 0666)

	return tempDir, err
}
