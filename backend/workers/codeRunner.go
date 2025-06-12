package workers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
)

// Each codeRunner instance gets a `job` and its own `containerManager`
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

// Creates a temporary file and writes the code in it
// The file is mounted in a Docker volume
// Docker reads that file and executes the code
func (cr *codeRunner) RunCode(outputChannel chan JobResult) {
	filename := LanguageMap[cr.job.Language].FileName
	tempDirPath, err := writeToTempFile(cr.job.Code, filename)

	if err != nil {
		outputChannel <- JobResult{Output: "", Error: err.Error()}
		return
	}

	defer os.RemoveAll(tempDirPath)

	imageName := LanguageMap[cr.job.Language].ImageName
	resp, err := cr.containerManager.CreateContainer(imageName, tempDirPath)
	if err != nil {
		outputChannel <- JobResult{Output: "", Error: err.Error()}
		return
	}

	err = cr.containerManager.StartContainer(resp.ID)

	if err != nil {
		outputChannel <- JobResult{Output: "", Error: err.Error()}
		return
	}

	// If the container exits early, send an error
	// Otherwise, send the result of the code
	// Remove the container in either case
	cr.containerManager.WaitForContainer(resp.ID, func(sc <-chan container.WaitResponse, ec <-chan error) {
		select {
		case err := <-ec:
			if err != nil {
				outputChannel <- JobResult{Output: "", Error: "Code took longer than 3 seconds to execute"}
				err = cr.containerManager.RemoveContainer(resp.ID)
				if err != nil {
					outputChannel <- JobResult{Output: "", Error: err.Error()}
				}
				close(outputChannel)
			}
		case <-sc:
			stdout, stderr, err := cr.containerManager.GetContainerOutputParsed(resp.ID)

			if err != nil {
				outputChannel <- JobResult{Output: "", Error: err.Error()}
				return
			}

			err = cr.containerManager.RemoveContainer(resp.ID)

			if err != nil {
				outputChannel <- JobResult{Output: "", Error: err.Error()}
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
			outputChannel <- JobResult{Output: output, Error: ""}
			close(outputChannel)
		}
	})
}

func writeToTempFile(code string, filename string) (string, error) {
	tempDir, err := os.MkdirTemp("", "temp")

	if err != nil {
		return "", err
	}

	err = os.WriteFile(filepath.Join(tempDir, filename), []byte(code), 0666)

	return tempDir, err
}
