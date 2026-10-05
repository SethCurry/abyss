package agentconfig

import "time"

// DefaultImage is the Docker image used when the agent config does not
// specify one.
const DefaultImage = "ghcr.io/sethcurry/abyss-pi:latest"

// DefaultServerPort is the container port the agent's API server listens on
// when the config does not specify one.
const DefaultServerPort = 8080

// Default paths where the agent's TLS certificates are stored inside the
// agent container. The client certificate and key are also stored there so
// persistent containers can be reconnected to in later sessions.
const (
	DefaultTLSServerCertPath = "/etc/abyss/tls/server.crt"
	DefaultTLSServerKeyPath  = "/etc/abyss/tls/server.key"
	DefaultTLSCACertPath     = "/etc/abyss/tls/ca.crt"
	DefaultTLSClientCertPath = "/etc/abyss/tls/client.crt"
	DefaultTLSClientKeyPath  = "/etc/abyss/tls/client.key"
)

// DefaultStartFilePath is the path inside the agent container whose
// existence signals that the agent has finished starting up.
const DefaultStartFilePath = "/tmp/abyss.start"

// WaitForStartFileSleepDuration is how long the client waits between checks
// for the start file to appear.
const WaitForStartFileSleepDuration = time.Second
