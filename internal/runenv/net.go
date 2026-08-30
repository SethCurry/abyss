package runenv

import (
	"fmt"
	"net"
)

// FreePort finds a random, unused TCP port on the host by asking the OS to
// bind to port 0 and reading back the assigned port. The listener is closed
// before returning, so the port is free for the caller to use.
func FreePort() (uint16, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("find free port: %w", err)
	}
	defer func() {
		_ = l.Close()
	}()

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("find free port: unexpected address type %T", l.Addr())
	}

	return uint16(addr.Port), nil
}
