package abyss

// ProxyLocation stores the location of a proxy,
// i.e. whether it is on the host or inside
// the container.
type ProxyLocation string

const (
	// LocationHost indicates the proxy is running on the host.
	LocationHost ProxyLocation = "host"

	// LocationContainer indicates the proxy is running in the container.
	LocationContainer ProxyLocation = "container"
)
