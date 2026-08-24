package runenv

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewContainer(docker *client.Client, containerID string) *Container {
	return &Container{
		client:      docker,
		logger:      log.With().Str("component", "container:"+containerID).Logger(),
		containerID: containerID,
	}
}

type Container struct {
	containerID string
	client      *client.Client
	logger      zerolog.Logger
}

// Stop stops and removes the container identified by containerID.
// It first sends a SIGTERM, waiting up to timeout for the container to exit,
// then falls back to SIGKILL if it is still running. The container is removed
// afterwards regardless of whether it exited on its own.
func (c *Container) Stop(ctx context.Context, timeout time.Duration) error {
	c.logger.Info().Msg("stopping container")

	if _, err := c.client.ContainerStop(ctx, c.containerID, client.ContainerStopOptions{Signal: "SIGTERM", Timeout: ptr(int(timeout.Seconds()))}); err != nil {
		c.logger.Warn().Err(err).Msg("container stop failed, attempting remove")
	}

	c.logger.Info().Msg("removing container")
	if _, err := c.client.ContainerRemove(ctx, c.containerID, client.ContainerRemoveOptions{Force: true}); err != nil {
		c.logger.Error().Err(err).Msg("failed to remove container")
		return fmt.Errorf("remove container: %w", err)
	}

	c.logger.Debug().Msg("container stopped and removed")
	return nil
}

// ExecBash runs script inside the container identified by containerID using
// "bash -c". It returns the script's stdout and stderr. If the script exits
// with a non-zero status, the returned error reports the exit code and the
// output is still returned so callers can inspect it.
func (c *Container) ExecBash(ctx context.Context, script string) (stdout, stderr string, err error) {
	c.logger.Debug().Msg("executing bash script in container")

	created, err := c.client.ExecCreate(ctx, c.containerID, client.ExecCreateOptions{
		Cmd:          []string{"bash", "-c", script},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to create exec")
		return "", "", fmt.Errorf("create exec: %w", err)
	}

	attach, err := c.client.ExecAttach(ctx, created.ID, client.ExecAttachOptions{})
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to attach to exec")
		return "", "", fmt.Errorf("attach to exec: %w", err)
	}
	defer attach.Close()

	var out, errOut bytes.Buffer
	if _, err := stdcopy.StdCopy(&out, &errOut, attach.Reader); err != nil {
		c.logger.Error().Err(err).Msg("failed to read exec output")
		return "", "", fmt.Errorf("read exec output: %w", err)
	}

	inspect, err := c.client.ExecInspect(ctx, created.ID, client.ExecInspectOptions{})
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to inspect exec")
		return "", "", fmt.Errorf("inspect exec: %w", err)
	}

	if inspect.ExitCode != 0 {
		return out.String(), errOut.String(), fmt.Errorf("script exited with code %d", inspect.ExitCode)
	}

	return out.String(), errOut.String(), nil
}

// CopyToContainer copies the file or directory at hostPath into the container
// identified by containerID, placing it under containerDir. Any parent
// directories of containerDir that do not already exist inside the container are
// created. The basename of hostPath is preserved, so a host path of
// "/tmp/foo.txt" copied into "/opt/app" lands at "/opt/app/foo.txt".
func (d *Container) CopyFromHost(ctx context.Context, containerID, hostPath, containerDir string) error {
	d.logger.Debug().
		Str("host_path", hostPath).
		Str("container_dir", containerDir).
		Msg("copying path into container")

	info, err := os.Stat(hostPath)
	if err != nil {
		return fmt.Errorf("stat host path %q: %w", hostPath, err)
	}

	// Ensure the destination directory exists inside the container before
	// extracting the tar archive into it.
	if _, _, err := d.ExecBash(ctx, fmt.Sprintf("mkdir -p %s", shellQuote(containerDir))); err != nil {
		return fmt.Errorf("create container directory %q: %w", containerDir, err)
	}

	pr, pw := io.Pipe()
	go func() {
		pw.CloseWithError(buildTar(pw, hostPath, info))
	}()

	if _, err := d.client.CopyToContainer(ctx, containerID, client.CopyToContainerOptions{
		DestinationPath: containerDir,
		Content:         pr,
	}); err != nil {
		d.logger.Error().Err(err).
			Str("host_path", hostPath).
			Str("container_dir", containerDir).
			Msg("failed to copy path into container")
		return fmt.Errorf("copy %q into container %q: %w", hostPath, containerDir, err)
	}

	d.logger.Debug().
		Str("host_path", hostPath).
		Str("container_dir", containerDir).
		Msg("copied path into container")
	return nil
}

// CopyFileFromHost copies a single file's content into the container
// identified by containerID, placing it at the exact containerPath. Parent
// directories of containerPath are created as needed. content is the file's
// bytes and mode sets the file permissions inside the container (only the
// permission bits are used).
func (d *Container) CopyFileFromHost(ctx context.Context, content []byte, containerPath string, mode os.FileMode) error {
	dir := filepath.Dir(containerPath)
	base := filepath.Base(containerPath)

	d.logger.Debug().
		Str("container_path", containerPath).
		Int("size", len(content)).
		Msg("copying file into container")

	if _, _, err := d.ExecBash(ctx, fmt.Sprintf("mkdir -p %s", shellQuote(dir))); err != nil {
		return fmt.Errorf("create container directory %q: %w", dir, err)
	}

	pr, pw := io.Pipe()
	go func() {
		pw.CloseWithError(writeContentTar(pw, content, base, mode))
	}()

	if _, err := d.client.CopyToContainer(ctx, d.containerID, client.CopyToContainerOptions{
		DestinationPath: dir,
		Content:         pr,
	}); err != nil {
		d.logger.Error().Err(err).
			Str("container_path", containerPath).
			Msg("failed to copy file into container")
		return fmt.Errorf("copy file into container %q: %w", containerPath, err)
	}

	d.logger.Debug().
		Str("container_path", containerPath).
		Msg("copied file into container")
	return nil
}
