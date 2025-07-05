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

func (g *GoldenFile) ListFiles(t *testing.T, namespace string) []string {
	t.Helper()
	
	// namespaceディレクトリのパスを構築
	goldenDir := filepath.Join("testdata", namespace, "golden")
	
	// パターンを構築（すべてのファイルを取得）
	pattern := filepath.Join(goldenDir, "*")
	
	// ファイル一覧を取得
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("failed to list files: %v", err)
	}
	
	// ディレクトリは除外してファイルのみを返す
	var result []string
	for _, file := range files {
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			result = append(result, file)
		}
	}
	
	return result
}

func GoldenFileWithPathHandler(f func(t *testing.T, namespace string, name string) string) GoldenFileOption {
	return func(g *GoldenFile) {
		g.pathHandler = f
	}
}
