//go:build windows

package acptools

import "os"

// exitStatus extracts the exit code from a finished process. On Windows the
// concept of a terminating signal does not map cleanly, so only the exit code
// is reported.
func exitStatus(ps *os.ProcessState, waitErr error) (exitCode *int, signal *string) {
	if ps == nil {
		return
	}
	code := ps.ExitCode()
	exitCode = &code
	return
}