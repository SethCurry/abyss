// Package timber contains logging utilities.
package timber

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/SethCurry/abyss/internal/fp"
)

const timeFormat = "2006-01-02_15-04-05"

func getLogDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	logDir := filepath.Join(home, ".local", "var", "abyss", "log")

	// create directories here since we return an error anyways
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to make parent directories of %q: %w", logDir, err)
	}

	return logDir, nil
}

// CleanLogDir removes old log entries to keep the directory clean.
// keepEntries specifies the number of entries to keep, with 0 meaning all entries.
func CleanLogDir(keepEntries int) error {
	logDir, err := getLogDir()
	if err != nil {
		return fmt.Errorf("failed to get log directory: %w", err)
	}

	files, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	files = fp.Filter(func(de os.DirEntry) bool {
		return strings.HasSuffix(de.Name(), ".log")
	}, files)

	if keepEntries == 0 {
		keepEntries = len(files)
	}

	if len(files) > keepEntries {
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})

		for i := 0; i < len(files)-keepEntries; i++ {
			if err := os.Remove(filepath.Join(logDir, files[i].Name())); err != nil {
				return fmt.Errorf("failed to remove log file %q: %w", files[i].Name(), err)
			}
		}
	}

	return nil
}

func OpenLogFile() (io.WriteCloser, error) {
	logDir, err := getLogDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get log directory: %w", err)
	}

	timestamp := time.Now().Format(timeFormat)

	logFilePath := filepath.Join(logDir, timestamp+".log")

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed open log file %q: %w", logFilePath, err)
	}

	return logFile, nil
}
