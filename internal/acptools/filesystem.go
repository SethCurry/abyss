package acptools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/coder/acp-go-sdk"
	"github.com/rs/zerolog"
)

func NewFilesystemTools(logger zerolog.Logger) *FilesystemTools {
	return &FilesystemTools{
		logger: logger,
	}
}

type FilesystemTools struct {
	logger zerolog.Logger
}

func (e *FilesystemTools) WriteTextFile(ctx context.Context, params acp.WriteTextFileRequest) (acp.WriteTextFileResponse, error) {
	e.logger.Debug().Str("method", "WriteTextFile").Msg("handling request")

	path := params.Path
	content := params.Content

	fd, err := os.Create(path)
	if err != nil {
		e.logger.Error().Err(err).Str("method", "WriteTextFile").Msg("failed to create file")
		return acp.WriteTextFileResponse{}, fmt.Errorf("failed to create or truncate %q: %w", path, err)
	}
	defer fd.Close()

	_, err = fd.Write([]byte(content))
	if err != nil {
		e.logger.Error().Err(err).Str("method", "WriteTextFile").Msg("failed to write to file")
		return acp.WriteTextFileResponse{}, fmt.Errorf("failed to write to %q: %w", path, err)
	}

	return acp.WriteTextFileResponse{}, nil
}

func (e *FilesystemTools) ReadTextFile(ctx context.Context, params acp.ReadTextFileRequest) (acp.ReadTextFileResponse, error) {
	var limit int
	if params.Limit != nil {
		limit = *params.Limit
	}

	var startLine int
	if params.Line != nil {
		startLine = *params.Line
	}

	currentLine := 0
	readLines := 0
	response := strings.Builder{}

	fd, err := os.Open(params.Path)
	if err != nil {
		e.logger.Error().Err(err).Str("method", "ReadTextFile").Msg("failed to open file")
		return acp.ReadTextFileResponse{}, fmt.Errorf("failed to open %q: %w", params.Path, err)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		currentLine += 1

		if currentLine >= startLine {
			response.WriteString(scanner.Text())
			readLines++
			if limit > 0 && readLines >= limit {
				break
			}
			response.WriteString("\n")
		}
	}

	if scanner.Err() != nil {
		e.logger.Error().Err(scanner.Err()).Str("method", "ReadTextFile").Msg("failed to scan file")
		return acp.ReadTextFileResponse{}, fmt.Errorf("failed to scan %q: %w", params.Path, scanner.Err())
	}

	return acp.ReadTextFileResponse{
		Content: response.String(),
	}, nil
}
