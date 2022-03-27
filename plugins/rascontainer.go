package plugins

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/net/context"
)

// StartContainer uses the Go SDK to run Docker containers..option 1
func StartConainer(imageWithTag string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

	reader, err := cli.ImagePull(ctx, imageWithTag, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	io.Copy(os.Stdout, reader)

	var chc *container.HostConfig
	var nnc *network.NetworkingConfig
	var vp *v1.Platform

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageWithTag,
		Cmd:   []string{"sleep", "5m"},
		Tty:   true,
	}, chc, nnc, vp, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, err
}

type ContainerParams struct {
	InputRasModelDir string `json:"input_ras_model_dir"`
	ModelName        string `json:"model_name"`
	PlanFile         string `json:"planfile"`
	OutputHDF        string `json:"output_hdf"`
	OutputLog        string `json:"output_log"`
}

func (cp ContainerParams) DirName() string {
	_, dirName := filepath.Split(cp.InputRasModelDir)
	return dirName
}

// RunSimInContainer uses system os to manage containers ...option 2
func RunSimInContainer(cp ContainerParams) (string, error) {

	rasContainer := "docker.io/lawlerseth/ras-docker-6.1-ubi8.5:latest"
	containerID, err := StartConainer(rasContainer)
	if err != nil {
		return "", err
	}

	// Wait for the container to boot and come online
	time.Sleep(5 * time.Second)

	// Director Mapping
	containerPath := fmt.Sprintf("%v:/sim", containerID)
	containerOutputPath := fmt.Sprintf("%v:/sim/%v/%v.tmp.hdf", containerID, cp.DirName(), cp.PlanFile)
	simDir := fmt.Sprintf("/sim/%v", cp.DirName())

	// System Docker calls
	copyModelInput := []string{"cp", cp.InputRasModelDir, containerPath}
	startSim := []string{"exec", containerID, "./run-model.sh", simDir, cp.ModelName}
	copyModelOutput := []string{"cp", containerOutputPath, cp.OutputHDF}
	// Todo: add simulation log mount or copy
	stopContainer := []string{"stop", containerID}
	removeContainer := []string{"rm", containerID}

	// Dump system calls into a slice for iteration loop
	systemCalls := [][]string{copyModelInput, startSim, copyModelOutput, stopContainer, removeContainer}

	for _, sysCall := range systemCalls {
		fmt.Println("Syscall: ", sysCall)
		cmd := exec.Command("docker", sysCall...)
		_, err = cmd.Output()
		if err != nil {
			// Comment next line to leave container running and debug via shell
			_ = exec.Command("docker", removeContainer...)
			fmt.Println("Terminating container due", containerID, err.Error())
			return "", err
		}

	}

	return "done", nil
}
