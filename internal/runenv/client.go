package runenv

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"net"
	"net/netip"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/moby/moby/api/types/container"

	"github.com/SethCurry/abyss/internal/agentconfig"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog"
)

// ContainerEndpoint describes how the host can reach a started container.
type ContainerEndpoint struct {
	// ContainerID is the ID of the started container.
	ContainerID string
	// IP is the address the host should connect to.
	IP string
	// Port is the host port bound to the container's exposed port.
	Port uint16
}

// String returns the endpoint as an "ip:port" address.
func (e ContainerEndpoint) String() string {
	return net.JoinHostPort(e.IP, strconv.FormatUint(uint64(e.Port), 10))
}

// DockerClient wraps the Docker SDK client and provides helpers for managing
// containers in a run environment.
type DockerClient struct {
	client *client.Client
	logger zerolog.Logger
}

// NewDockerClient creates a DockerClient configured from the environment
// (DOCKER_HOST, DOCKER_API_VERSION, DOCKER_CERT_PATH, DOCKER_TLS_VERIFY).
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}

	return &DockerClient{
		client: cli,
		logger: zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
	}, nil
}

// Close releases resources held by the underlying Docker client.
func (d *DockerClient) Close() error {
	d.logger.Debug().Msg("closing docker client")
	if err := d.client.Close(); err != nil {
		d.logger.Error().Err(err).Msg("failed to close docker client")
		return err
	}
	return nil
}

func (d *DockerClient) AbyssContainers(ctx context.Context) ([]container.Summary, error) {
	resp, err := d.client.ContainerList(ctx, client.ContainerListOptions{
		Filters: client.Filters{
			"label": map[string]bool{
				"abyss": true,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	return resp.Items, nil
}

// StartContainer pulls imageRef (if necessary) and starts a container from it.
// The container's containerPort is published to the host so the host can reach
// it. hostPort selects the host port to bind; pass 0 to let Docker assign a
// free port. config and hostConfig may be nil to use Docker's defaults for the
// image.
func (d *DockerClient) StartContainer(
	ctx context.Context,
	imageRef string,
	config *container.Config,
	hostConfig *container.HostConfig,
	name string,
	containerPort uint16,
	hostPort uint16,
) (*Container, ContainerEndpoint, error) {
	d.logger.Debug().
		Str("image", imageRef).
		Str("name", name).
		Uint16("container_port", containerPort).
		Uint16("host_port", hostPort).
		Msg("starting container")

	//if err := d.pullImage(ctx, imageRef); err != nil {
	//	return ContainerEndpoint{}, err
	//}

	port, ok := network.PortFrom(containerPort, network.TCP)
	if !ok {
		return nil, ContainerEndpoint{}, fmt.Errorf("invalid container port %d", containerPort)
	}

	if config == nil {
		config = &container.Config{}
	}
	config.Image = imageRef
	config.ExposedPorts = network.PortSet{port: {}}
	config.Env = []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}
	config.Labels = map[string]string{
		"abyss": "true",
	}

	if hostConfig == nil {
		hostConfig = &container.HostConfig{}
	}
	hostConfig.PortBindings = network.PortMap{
		port: {{
			HostIP:   netip.IPv4Unspecified(),
			HostPort: strconv.FormatUint(uint64(hostPort), 10),
		}},
	}

	created, err := d.client.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:     config,
		HostConfig: hostConfig,
		Name:       name,
	})
	if err != nil {
		d.logger.Error().Err(err).Str("image", imageRef).Str("name", name).Msg("failed to create container")
		return nil, ContainerEndpoint{}, fmt.Errorf("create container: %w", err)
	}

	if _, err := d.client.ContainerStart(ctx, created.ID, client.ContainerStartOptions{}); err != nil {
		d.logger.Error().Err(err).Str("container_id", created.ID).Msg("failed to start container")
		return nil, ContainerEndpoint{}, fmt.Errorf("start container: %w", err)
	}

	// Resolve the actual host port in case Docker assigned one for us.
	actualPort := hostPort
	if hostPort == 0 {
		inspected, err := d.client.ContainerInspect(ctx, created.ID, client.ContainerInspectOptions{})
		if err != nil {
			d.logger.Error().Err(err).Str("container_id", created.ID).Msg("failed to inspect container")
			return nil, ContainerEndpoint{}, fmt.Errorf("inspect container: %w", err)
		}

		bindings := inspected.Container.NetworkSettings.Ports[port]
		if len(bindings) == 0 {
			return nil, ContainerEndpoint{}, fmt.Errorf("no host port binding found for container port %d", containerPort)
		}

		p, err := strconv.ParseUint(bindings[0].HostPort, 10, 16)
		if err != nil {
			return nil, ContainerEndpoint{}, fmt.Errorf("parse host port %q: %w", bindings[0].HostPort, err)
		}
		actualPort = uint16(p)
	}

	endpoint := ContainerEndpoint{
		ContainerID: created.ID,
		IP:          d.hostIP(),
		Port:        actualPort,
	}
	d.logger.Debug().
		Str("container_id", endpoint.ContainerID).
		Str("endpoint", endpoint.String()).
		Msg("container started")

	return NewContainer(d, created.ID), endpoint, nil
}

func cleanPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory to resolve ~: %w", err)
		}

		path = strings.Replace(path, "~", home, 1)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %q: %w", path, err)
	}

	return filepath.Clean(absPath), nil
}

// ApplyHostMounts populates hostConfig.Binds from the HostMounts declared in
// cfg. Each entry maps a host path to a container path using Docker's
// "hostPath:containerPath" bind format. If hostConfig is nil a new HostConfig is
// allocated. Mounts are appended to any existing binds so callers can layer
// additional mounts. Paths are normalized to absolute form; a mount whose host
// path does not exist on the host is skipped with a warning, since Docker would
// create it as an empty directory owned by root, which is rarely intended.
func (d *DockerClient) ApplyHostMounts(cfg *agentconfig.DockerConfig, hostConfig *container.HostConfig) *container.HostConfig {
	if hostConfig == nil {
		hostConfig = &container.HostConfig{}
	}
	if cfg == nil || len(cfg.HostMounts) == 0 {
		return hostConfig
	}

	for _, mount := range cfg.HostMounts {
		hostPath, err := cleanPath(mount.Source)
		if err != nil {
			d.logger.Err(err).Msg("failed to get absolute host path")
			return nil
		}

		var containerPath string

		if mount.Destination != "" {
			containerPath, err = cleanPath(mount.Destination)
			if err != nil {
				d.logger.Err(err).Msg("failed to get absolute container path")
				return nil
			}
		} else {
			containerPath = hostPath
		}

		if !filepath.IsAbs(hostPath) {
			d.logger.Warn().Str("host_path", hostPath).Str("container_path", containerPath).Msg("skipping non-absolute host mount")
			continue
		}
		if !filepath.IsAbs(containerPath) {
			d.logger.Warn().Str("host_path", hostPath).Str("container_path", containerPath).Msg("skipping mount with non-absolute container path")
			continue
		}
		if _, err := os.Stat(hostPath); err != nil {
			d.logger.Warn().Err(err).Str("host_path", hostPath).Str("container_path", containerPath).Msg("skipping host mount, path does not exist on host")
			continue
		}

		d.logger.Debug().Str("host_path", hostPath).Str("container_path", containerPath).Msg("adding host mount")
		hostConfig.Binds = append(hostConfig.Binds, hostPath+":"+containerPath)
	}

	return hostConfig
}

// ptr returns a pointer to v.
func ptr[T any](v T) *T { return &v }

// writeContentTar writes a tar archive containing a single regular file entry
// named name with the given content and permission bits from mode.
func writeContentTar(w io.Writer, content []byte, name string, mode os.FileMode) error {
	tw := tar.NewWriter(w)
	hdr := &tar.Header{
		Name: name,
		Mode: int64(mode.Perm()),
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("write tar header for %q: %w", name, err)
	}
	if _, err := tw.Write(content); err != nil {
		return fmt.Errorf("write tar content for %q: %w", name, err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("close tar writer: %w", err)
	}
	return nil
}

// buildTar writes a tar archive of src to w. Entries are named relative to src
// and prefixed with the basename of src, so the basename of src becomes the
// top-level entry in the archive (e.g. /tmp/foo/sub/a.txt becomes "foo/sub/a.txt").
// This preserves the nested directory structure when the archive is extracted
// into a destination directory via the Docker archive API, matching the
// semantics of `docker cp`.
func buildTar(w io.Writer, src string, info os.FileInfo) error {
	tw := tar.NewWriter(w)

	closeErr := func() error {
		if err := tw.Close(); err != nil {
			return fmt.Errorf("close tar writer: %w", err)
		}
		return nil
	}

	if !info.IsDir() {
		if err := writeTarEntry(tw, src, filepath.Base(src), info); err != nil {
			return err
		}
		return closeErr()
	}

	base := filepath.Base(src)
	err := filepath.Walk(src, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("rel path %q: %w", path, err)
		}
		rel = filepath.ToSlash(rel)

		// Prefix every entry with the basename of src so that src is preserved
		// as the top-level directory in the archive. Without this the contents
		// of src are spread flat across the extraction destination and the
		// nested directory structure rooted at src is lost.
		var name string
		if rel == "." {
			name = base
		} else {
			name = base + "/" + rel
		}
		if fi.IsDir() {
			name += "/"
		}
		return writeTarEntry(tw, path, name, fi)
	})
	if err != nil {
		return err
	}
	return closeErr()
}

// writeTarEntry writes a single entry to tw for the host path src, using name as
// the entry's name inside the archive. Directories write only a header;
// symlinks write a header with the link target; regular files copy their
// contents.
func writeTarEntry(tw *tar.Writer, src, name string, fi os.FileInfo) error {
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return fmt.Errorf("build tar header for %q: %w", src, err)
	}
	hdr.Name = name

	if fi.Mode()&os.ModeSymlink != 0 {
		link, err := os.Readlink(src)
		if err != nil {
			return fmt.Errorf("read symlink %q: %w", src, err)
		}
		hdr.Linkname = link
		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("write tar header for %q: %w", src, err)
		}
		return nil
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("write tar header for %q: %w", src, err)
	}
	if !fi.Mode().IsRegular() {
		return nil
	}

	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %q: %w", src, err)
	}
	defer f.Close()
	if _, err := io.Copy(tw, f); err != nil {
		return fmt.Errorf("copy %q into tar: %w", src, err)
	}
	return nil
}

// shellQuote wraps s in single quotes, escaping any embedded single quotes so
// that the result can be used safely as a single shell argument.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}

// hostIP returns the address the host should use to reach the Docker daemon.
// For a local daemon (unix socket or named pipe) this is the loopback address;
// for a remote daemon it is the daemon's host.
func (d *DockerClient) hostIP() string {
	u, err := url.Parse(d.client.DaemonHost())
	if err != nil {
		d.logger.Warn().Err(err).Str("daemon_host", d.client.DaemonHost()).Msg("failed to parse docker daemon host, defaulting to loopback")
		return "127.0.0.1"
	}

	switch u.Scheme {
	case "tcp", "http", "https":
		if host := u.Hostname(); host != "" {
			return host
		}
	}

	return "127.0.0.1"
}

// pullImage ensures imageRef is present on the Docker host.
func (d *DockerClient) pullImage(ctx context.Context, imageRef string) error {
	d.logger.Debug().Str("image", imageRef).Msg("pulling image")

	resp, err := d.client.ImagePull(ctx, imageRef, client.ImagePullOptions{})
	if err != nil {
		d.logger.Error().Err(err).Str("image", imageRef).Msg("failed to pull image")
		return fmt.Errorf("pull image %q: %w", imageRef, err)
	}
	defer resp.Close()

	if err := resp.Wait(ctx); err != nil {
		d.logger.Error().Err(err).Str("image", imageRef).Msg("failed to pull image")
		return fmt.Errorf("pull image %q: %w", imageRef, err)
	}

	return nil
}
