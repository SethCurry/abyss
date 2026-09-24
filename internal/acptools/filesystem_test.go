package acptools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/coder/acp-go-sdk"
	"github.com/rs/zerolog"
)

func newTestFilesystemTools() *FilesystemTools {
	return NewFilesystemTools(zerolog.Nop())
}

// helper: write a file with the given content using the tool under test.
func writeFile(t *testing.T, ft *FilesystemTools, path, content string) {
	t.Helper()
	if _, err := ft.WriteTextFile(context.Background(), acp.WriteTextFileRequest{
		Path:    path,
		Content: content,
	}); err != nil {
		t.Fatalf("WriteTextFile(%q): %v", path, err)
	}
}

func TestWriteTextFileCreatesFile(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "hello.txt")

	writeFile(t, ft, path, "hello world")

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile: %v", err)
	}
	if string(got) != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", string(got))
	}
}

func TestWriteTextFileOverwritesExisting(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, ft, path, "first")
	writeFile(t, ft, path, "second")

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile: %v", err)
	}
	if string(got) != "second" {
		t.Fatalf("expected overwritten content %q, got %q", "second", string(got))
	}
}

func TestWriteTextFileTruncates(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, ft, path, "a much longer original content")
	writeFile(t, ft, path, "short")

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile: %v", err)
	}
	if string(got) != "short" {
		t.Fatalf("expected truncated content %q, got %q", "short", string(got))
	}
}

func TestWriteTextFileEmptyContent(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	writeFile(t, ft, path, "")

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("os.Stat: %v", err)
	}
	if info.Size() != 0 {
		t.Fatalf("expected empty file, got size %d", info.Size())
	}
}

func TestWriteTextFileMissingDirectory(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "does", "not", "exist", "file.txt")

	_, err := ft.WriteTextFile(context.Background(), acp.WriteTextFileRequest{
		Path:    path,
		Content: "data",
	})
	if err == nil {
		t.Fatal("expected error writing to missing directory")
	}
	if _, statErr := os.Stat(path); statErr == nil {
		t.Fatal("expected file to not be created")
	}
}

func TestReadTextFileFull(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	content := "line1\nline2\nline3\n"
	writeFile(t, ft, path, content)

	resp, err := ft.ReadTextFile(context.Background(), acp.ReadTextFileRequest{
		Path: path,
	})
	if err != nil {
		t.Fatalf("ReadTextFile: %v", err)
	}
	// Each scanned line gets a trailing newline appended, so the output
	// matches the original content (which also ends with a newline).
	if resp.Content != content {
		t.Fatalf("expected %q, got %q", content, resp.Content)
	}
}

func TestReadTextFileNoTrailingNewline(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, ft, path, "line1\nline2\nline3")

	resp, err := ft.ReadTextFile(context.Background(), acp.ReadTextFileRequest{
		Path: path,
	})
	if err != nil {
		t.Fatalf("ReadTextFile: %v", err)
	}
	// The scanner skips the trailing newline but the tool re-adds one per
	// line, so the result is normalized to end with a newline.
	if resp.Content != "line1\nline2\nline3\n" {
		t.Fatalf("expected %q, got %q", "line1\nline2\nline3\n", resp.Content)
	}
}

func TestReadTextFileEmptyFile(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	writeFile(t, ft, path, "")

	resp, err := ft.ReadTextFile(context.Background(), acp.ReadTextFileRequest{
		Path: path,
	})
	if err != nil {
		t.Fatalf("ReadTextFile: %v", err)
	}
	if resp.Content != "" {
		t.Fatalf("expected empty content, got %q", resp.Content)
	}
}

func TestReadTextFileWithLineOffset(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, ft, path, "line1\nline2\nline3\nline4\n")

	line := 3
	resp, err := ft.ReadTextFile(context.Background(), acp.ReadTextFileRequest{
		Path: path,
		Line: &line,
	})
	if err != nil {
		t.Fatalf("ReadTextFile: %v", err)
	}
	// Line is 1-indexed; starting at line 3 yields lines 3 and 4.
	if resp.Content != "line3\nline4\n" {
		t.Fatalf("expected %q, got %q", "line3\nline4\n", resp.Content)
	}
}

func TestReadTextFileWithLimit(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, ft, path, "line1\nline2\nline3\nline4\nline5\n")

	limit := 2
	resp, err := ft.ReadTextFile(context.Background(), acp.ReadTextFileRequest{
		Path:  path,
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("ReadTextFile: %v", err)
	}
	// The limit break fires before appending the trailing newline, so the
	// last included line has no newline.
	if resp.Content != "line1\nline2" {
		t.Fatalf("expected %q, got %q", "line1\nline2", resp.Content)
	}
}

func TestReadTextFileWithLineOffsetAndLimit(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, ft, path, "line1\nline2\nline3\nline4\nline5\n")

	line := 2
	limit := 2
	resp, err := ft.ReadTextFile(context.Background(), acp.ReadTextFileRequest{
		Path:  path,
		Line:  &line,
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("ReadTextFile: %v", err)
	}
	// Start at line 2, take 2 lines -> lines 2 and 3 (last line no newline).
	if resp.Content != "line2\nline3" {
		t.Fatalf("expected %q, got %q", "line2\nline3", resp.Content)
	}
}

func TestReadTextFileLineOffsetBeyondEOF(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, ft, path, "line1\nline2\n")

	line := 10
	resp, err := ft.ReadTextFile(context.Background(), acp.ReadTextFileRequest{
		Path: path,
		Line: &line,
	})
	if err != nil {
		t.Fatalf("ReadTextFile: %v", err)
	}
	if resp.Content != "" {
		t.Fatalf("expected empty content past EOF, got %q", resp.Content)
	}
}

func TestReadTextFileMissingFile(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "nope.txt")

	_, err := ft.ReadTextFile(context.Background(), acp.ReadTextFileRequest{
		Path: path,
	})
	if err == nil {
		t.Fatal("expected error reading missing file")
	}
	if !strings.Contains(err.Error(), "failed to open") {
		t.Fatalf("expected error about opening file, got %v", err)
	}
}

func TestReadTextFileLimitZeroMeansAll(t *testing.T) {
	ft := newTestFilesystemTools()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	content := "a\nb\nc\nd\n"
	writeFile(t, ft, path, content)

	limit := 0
	resp, err := ft.ReadTextFile(context.Background(), acp.ReadTextFileRequest{
		Path:  path,
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("ReadTextFile: %v", err)
	}
	// A limit of 0 means "no limit", so all lines are returned.
	if resp.Content != content {
		t.Fatalf("expected %q, got %q", content, resp.Content)
	}
}