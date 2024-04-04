package template

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type DockerImageWriter struct{}

type ProgressDetail struct {
	Current uint64 `json:"current"`
	Total   uint64 `json:"total"`
}

type DockerStatusMessage struct {
	Status         string         `json:"status"`
	ID             string         `json:"id"`
	ProgressDetail ProgressDetail `json:"progressDetail"`
	Progress       string         `json:"progress"`
}

func (w *DockerImageWriter) Write(p []byte) (n int, err error) {
	// fmt.Println(string(p))

	msg := DockerStatusMessage{}
	if err := json.Unmarshal(p, &msg); err != nil {
		return 0, err
	}

	var percentComplete uint8
	if msg.ProgressDetail.Total > 0 {
		percentComplete = uint8((float64(msg.ProgressDetail.Current) / float64(msg.ProgressDetail.Total)) * 100)
	}
	if strings.Compare(msg.Status, "Download complete") == 0 {
		percentComplete = 100
	}
	// fmt.Printf("%+v\n", msg)
	fmt.Printf("%d percent complete!\n", percentComplete)

	return len(p), nil
}

func (t *Template) executeActionStep(step Step, _ afero.Fs, outputFs afero.BasePathFs) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	fmt.Printf("%+v\n", cli)

	reader, err := cli.ImagePull(ctx, step.Action, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return errors.Wrap(err, "error pull output to stdout")
	}

	outputRealPath, err := outputFs.RealPath(".")
	if err != nil {
		return err
	}
	outputSource, err := filepath.Abs(outputRealPath)
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      step.Action,
		Tty:        false,
		Cmd:        step.Command,
		WorkingDir: "/project",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:     "bind",
				Source:   outputSource,
				Target:   "/project",
				ReadOnly: false,
			},
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return nil
}
