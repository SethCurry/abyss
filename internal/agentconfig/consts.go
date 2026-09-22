package agentconfig

import "time"

// DefaultImage is the Docker image used when the agent config does not
// specify one.
const DefaultImage = "ghcr.io/sethcurry/abyss-pi:latest"

// DefaultServerPort is the container port the agent's API server listens on
// when the config does not specify one.
const DefaultServerPort = 8080

// Default paths where the agent's TLS server certificate, key, and CA
// certificate are mounted inside the agent container.
const (
	DefaultTLSServerCertPath = "/etc/abyss/tls/server.crt"
	DefaultTLSServerKeyPath  = "/etc/abyss/tls/server.key"
	DefaultTLSCACertPath     = "/etc/abyss/tls/ca.crt"
)

// DefaultStartFilePath is the path inside the agent container whose
// existence signals that the agent has finished starting up.
const DefaultStartFilePath = "/tmp/abyss.start"

// WaitForStartFileSleepDuration is how long the client waits between checks
// for the start file to appear.
const WaitForStartFileSleepDuration = time.Second
