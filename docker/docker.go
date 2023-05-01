package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Docker struct {
	cli *client.Client
}

func New() (docker *Docker, err error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return
	}
	docker = &Docker{cli: cli}
	return
}

// StopAllContainers - stops all running docker containers
func (d *Docker) StopAllContainers(ctx context.Context) (err error) {
	list, err := d.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range list {
		err = d.cli.ContainerStop(ctx, container.ID, nil)
		if err != nil {
			return
		}
	}
	return
}

type ContainerOptions struct {
	Image       string
	Name        string
	Host        string
	InsidePort  int
	ExposedPort int
	Envs        map[string]string
}

// DownloadImage - downloads docker image
func (d *Docker) DownloadImage(ctx context.Context, imageName string) (err error) {
	reader, err := d.cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return
	}
	io.Copy(os.Stdout, reader)
	return
}

// DeleteContainer - deletes container
func (d *Docker) DeleteContainer(ctx context.Context, containerName string) (err error) {
	list, err := d.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return
	}
	for _, container := range list {
		if container.Names[0] == "/"+containerName {
			d.cli.ContainerStop(ctx, container.ID, nil)
			err = d.cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
			if err != nil {
				return err
			}
		}
	}
	return
}

// CreateContainer - creates docker container and runs it
func (d *Docker) CreateContainer(ctx context.Context, opts ContainerOptions) (containerID string, err error) {
	envs := []string{}
	for key, val := range opts.Envs {
		envs = append(envs, fmt.Sprintf("%s=%s", key, val))
	}

	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image:        opts.Image,
		ExposedPorts: nat.PortSet{nat.Port(strconv.Itoa(opts.ExposedPort)): struct{}{}},
		Env:          envs,
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{nat.Port(strconv.Itoa(opts.ExposedPort)): {{HostIP: opts.Host, HostPort: strconv.Itoa(opts.InsidePort)}}},
	}, nil, nil, opts.Name)
	if err != nil {
		return
	}
	containerID = resp.ID

	if err := d.cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return containerID, err
	}
	return
}

func (d *Docker) IsContainerRunning(ctx context.Context, containerID string) (running bool, err error) {
	attempts := 3
	for attempts > 0 {
		info, err := d.cli.ContainerInspect(ctx, containerID)
		if err == nil && info.State.Running {
			running = true
			break
		}
		time.Sleep(time.Duration(time.Second))
		attempts--
	}
	return
}
