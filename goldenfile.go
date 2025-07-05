package gilt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var _ interface {
	ReadWriter
	Opener
	Lister
} = (*GoldenFile)(nil)

// GoldenFile implements ReadWriter, Opener, and Lister interfaces for file operations.
type GoldenFile struct {
	namespace   string
	pathHandler func(t *testing.T, namespace string, name string) string
	tracker     *FileTracker
}

// GoldenFileOption configures a GoldenFile instance.
type GoldenFileOption func(*GoldenFile)

// NewGoldenFile creates a new GoldenFile instance with the specified namespace and options.
func NewGoldenFile(namespace string, opts ...GoldenFileOption) *GoldenFile {
	f := &GoldenFile{
		namespace:   namespace,
		pathHandler: DefaultPathHandler,
		tracker:     NewFileTracker(),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// DefaultPathHandler generates the default file path for golden files.
func DefaultPathHandler(t *testing.T, namespace string, name string) string {
	t.Helper()
	return filepath.Join("testdata", namespace, "golden", fmt.Sprintf("%s.golden", name))
}

func (g *GoldenFile) filepath(t *testing.T, name string) string {
	t.Helper()
	return g.pathHandler(t, g.namespace, name)
}

// Read reads the contents of a golden file and tracks its usage.
func (g *GoldenFile) Read(t *testing.T, name string) []byte {
	t.Helper()
	filePath := g.filepath(t, name)

	// Track file usage
	if g.tracker != nil {
		g.tracker.TrackUsage(t, filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// Write writes data to a golden file and tracks its usage.
func (g *GoldenFile) Write(t *testing.T, name string, data []byte) {
	t.Helper()
	filePath := g.filepath(t, name)

	// Track file usage
	if g.tracker != nil {
		g.tracker.TrackUsage(t, filePath)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		t.Fatal(err)
	}

	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

// Open opens a golden file and returns a file handle.
func (g *GoldenFile) Open(t *testing.T, name string) *os.File {
	t.Helper()
	filePath := g.filepath(t, name)
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

// ListFiles returns a list of all golden files in the specified namespace.
func (g *GoldenFile) ListFiles(t *testing.T, namespace string) []string {
	t.Helper()

	// Build path to namespace directory
	goldenDir := filepath.Join("testdata", namespace, "golden")

	// Build pattern to get all files
	pattern := filepath.Join(goldenDir, "*")

	// Get file list
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("failed to list files: %v", err)
	}

	// Exclude directories and return only files
	var result []string
	for _, file := range files {
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			result = append(result, file)
		}
	}

	return result
}

// GetUnusedFiles returns a list of files that have not been accessed during the test.
func (g *GoldenFile) GetUnusedFiles(t *testing.T, namespace string) []string {
	t.Helper()
	allFiles := g.ListFiles(t, namespace)
	if g.tracker != nil {
		return g.tracker.GetUnusedFiles(allFiles)
	}
	return []string{}
}

// GoldenFileWithPathHandler configures a GoldenFile to use a custom path handler.
func GoldenFileWithPathHandler(f func(t *testing.T, namespace string, name string) string) GoldenFileOption {
	return func(g *GoldenFile) {
		g.pathHandler = f
	}
}

// GoldenFileWithTracker configures a GoldenFile to use a specific FileTracker instance.
func GoldenFileWithTracker(tracker *FileTracker) GoldenFileOption {
	return func(g *GoldenFile) {
		g.tracker = tracker
	}
}
