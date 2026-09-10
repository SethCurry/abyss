package agentconfig

import "time"

const (
	DefaultImage      = "ghcr.io/sethcurry/abyss-pi:latest"
	DefaultServerPort = 8080

	DefaultTLSServerCertPath      = "/etc/abyss/tls/server.crt"
	DefaultTLSServerKeyPath       = "/etc/abyss/tls/server.key"
	DefaultTLSCACertPath          = "/etc/abyss/tls/ca.crt"
	DefaultStartFilePath          = "/tmp/abyss.start"
	WaitForStartFileSleepDuration = time.Second
)
