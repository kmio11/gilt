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
} = (*GoldenFile)(nil)

type GoldenFile struct {
	namespace   string
	pathHandler func(t *testing.T, namespace string, name string) string
}

type GoldenFileOption func(*GoldenFile)

func NewGoldenFile(namespace string, opts ...GoldenFileOption) *GoldenFile {
	f := &GoldenFile{
		namespace:   namespace,
		pathHandler: DefaultPathHandler,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func DefaultPathHandler(t *testing.T, namespace string, name string) string {
	t.Helper()
	return filepath.Join("testdata", namespace, "golden", fmt.Sprintf("%s.golden", name))
}

func (g *GoldenFile) filepath(t *testing.T, name string) string {
	t.Helper()
	return g.pathHandler(t, g.namespace, name)
}

func (g *GoldenFile) Read(t *testing.T, name string) []byte {
	t.Helper()
	filePath := g.filepath(t, name)
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func (g *GoldenFile) Write(t *testing.T, name string, data []byte) {
	t.Helper()
	filePath := g.filepath(t, name)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		t.Fatal(err)
	}

	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func (g *GoldenFile) Open(t *testing.T, name string) *os.File {
	t.Helper()
	filePath := g.filepath(t, name)
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func GoldenFileWithPathHandler(f func(t *testing.T, namespace string, name string) string) GoldenFileOption {
	return func(g *GoldenFile) {
		g.pathHandler = f
	}
}
