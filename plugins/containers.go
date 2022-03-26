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
		Cmd:   []string{"sleep", "2m"},
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

	containerPath := fmt.Sprintf("%v:/sim", containerID)
	fmt.Println(containerID, containerPath)

	cmd := exec.Command("docker", "cp", cp.InputRasModelDir, containerPath)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	simDir := fmt.Sprintf("/sim/%v", cp.DirName())

	cmd = exec.Command("docker", "exec", containerID, "./run-model.sh", simDir, cp.ModelName)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	containerOutputPath := fmt.Sprintf("%v:/sim/%v/%v.tmp.hdf", containerID, cp.DirName(), cp.PlanFile)

	cmd = exec.Command("docker", "cp", containerOutputPath, cp.OutputHDF)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(stdout, err.Error())
		return "", err
	}

	cmd = exec.Command("docker", "stop", containerID)

	stdout, err = cmd.Output()
	if err != nil {
		fmt.Println(stdout, err.Error())
		return "", err
	}

	cmd = exec.Command("docker", "rm", containerID)

	stdout, err = cmd.Output()
	if err != nil {
		fmt.Println(stdout, err.Error())
		return "", err
	}

	return "done", nil
}
