package workers

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type WaitContainerCallback = func(<-chan container.WaitResponse, <-chan error)

type containerManager struct {
	cli *client.Client
}

func CreateContainerManager() (*containerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}

	return &containerManager{cli}, nil
}

func (manager *containerManager) CreateContainer(imageName, volumePath string) (container.CreateResponse, error) {
	resp, err := manager.cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: imageName,
		},
		&container.HostConfig{
			Binds:      []string{volumePath + ":/app"},
			AutoRemove: true,
		}, nil, nil, "")

	if err != nil {
		return container.CreateResponse{}, err
	}

	return resp, nil
}

func (manager *containerManager) StartContainer(containerId string) error {
	return manager.cli.ContainerStart(context.Background(), containerId, container.StartOptions{})
}

func (manager *containerManager) WaitForContainer(containerId string, callbackFn WaitContainerCallback) {
	statusCh, errCh := manager.cli.ContainerWait(context.Background(), containerId, container.WaitConditionNotRunning)
	callbackFn(statusCh, errCh)
}

func (manager *containerManager) GetContainerOutput(containerId string) (logs io.Reader, err error) {
	rawLogs, err := manager.cli.ContainerLogs(context.Background(), containerId, container.LogsOptions{ShowStdout: true, ShowStderr: true})

	if err != nil {
		return nil, err
	}

	return rawLogs, err
}

func (manager *containerManager) GetContainerOutputParsed(containerId string) (stdout, stderr string, err error) {
	rawLogs, err := manager.cli.ContainerLogs(context.Background(), containerId, container.LogsOptions{ShowStdout: true, ShowStderr: true})

	if err != nil {
		return "", "", err
	}
	defer rawLogs.Close()

	var stdoutBuf, stderrBuf io.Writer
	stdoutBuilder := &strings.Builder{}
	stderrBuilder := &strings.Builder{}

	stdoutBuf = stdoutBuilder
	stderrBuf = stderrBuilder

	_, err = stdcopy.StdCopy(stdoutBuf, stderrBuf, rawLogs)
	if err != nil {
		return "", "", err
	}

	return stdoutBuilder.String(), stderrBuilder.String(), nil
}

func (manager *containerManager) RemoveContainer(containerId string) error {
	return manager.cli.ContainerRemove(context.Background(), containerId, container.RemoveOptions{})
}
