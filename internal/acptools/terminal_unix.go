//go:build !windows

package acptools

import (
	"os"
	"syscall"
)

// exitStatus extracts the exit code and terminating signal (if any) from a
// finished process. On Unix-like systems this uses the wait status reported by
// the kernel.
func exitStatus(ps *os.ProcessState, waitErr error) (exitCode *int, signal *string) {
	if ps == nil {
		return
	}
	if ws, ok := ps.Sys().(syscall.WaitStatus); ok {
		if ws.Signaled() {
			s := ws.Signal().String()
			signal = &s
			return
		}
		if ws.Exited() {
			code := ws.ExitStatus()
			exitCode = &code
			return
		}
	}
	// Fallback for unexpected process states.
	code := ps.ExitCode()
	exitCode = &code
	return
}