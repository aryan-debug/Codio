package workers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
)

type codeRunner struct {
	job              Job
	containerManager *containerManager
}

func CreateCodeRunner(job Job) (*codeRunner, error) {
	containerManager, err := CreateContainerManager()

	if err != nil {
		return nil, err
	}

	return &codeRunner{
		job:              job,
		containerManager: containerManager,
	}, nil
}

func (cr *codeRunner) RunCode(outputChannel chan JobResult) {
	filename := "test.py"
	if cr.job.Language == Java {
		filename = "Main.java"
	}

	tempDirPath, err := writeToTempFile(cr.job.Code, filename)

	if err != nil {
		outputChannel <- JobResult{Output: "", Error: err}
		return
	}

	defer os.RemoveAll(tempDirPath)

	imageName := LanguageImageMap[cr.job.Language]

	resp, err := cr.containerManager.CreateContainer(imageName, tempDirPath)

	if err != nil {
		outputChannel <- JobResult{Output: "", Error: err}
		return
	}

	err = cr.containerManager.StartContainer(resp.ID)

	if err != nil {
		outputChannel <- JobResult{Output: "", Error: err}
		return
	}

	cr.containerManager.WaitForContainer(resp.ID, func(sc <-chan container.WaitResponse, ec <-chan error) {
		<-sc
	})

	stdout, stderr, err := cr.containerManager.GetContainerOutputParsed(resp.ID)

	if err != nil {
		outputChannel <- JobResult{Output: "", Error: err}
		return
	}

	err = cr.containerManager.RemoveContainer(resp.ID)

	if err != nil {
		outputChannel <- JobResult{Output: "", Error: err}
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
	fmt.Println(output)
	outputChannel <- JobResult{Output: output, Error: nil}
}

func writeToTempFile(code string, filename string) (string, error) {
	tempDir, err := os.MkdirTemp("", "temp")

	if err != nil {
		return "", err
	}

	err = os.WriteFile(filepath.Join(tempDir, filename), []byte(code), 0666)

	return tempDir, err
}
