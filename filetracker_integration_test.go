package gilt_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kmio11/gilt"
)

func TestFileTracker_Integration(t *testing.T) {
	namespace := "TestFileTracker_Integration"

	// テスト用のディレクトリとファイルを作成
	testDir := filepath.Join("testdata", namespace, "golden")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// クリーンアップ
	t.Cleanup(func() {
		os.RemoveAll(filepath.Join("testdata", namespace))
	})

	// テスト用のファイルを作成
	testFiles := []string{"used.golden", "unused.golden"}
	for _, filename := range testFiles {
		filePath := filepath.Join(testDir, filename)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// GoldenFileを作成（内蔵のFileTrackerを使用）
	goldenFile := gilt.NewGoldenFile(namespace)

	// ファイルを読み込み（自動的に追跡される）
	data := goldenFile.Read(t, "used")
	if string(data) != "test content" {
		t.Errorf("expected 'test content', got '%s'", string(data))
	}

	// 未使用ファイルを取得
	unusedFiles := goldenFile.GetUnusedFiles(t, namespace)

	// 検証
	if len(unusedFiles) != 1 {
		t.Errorf("expected 1 unused file, got %d", len(unusedFiles))
	}
	if len(unusedFiles) > 0 && filepath.Base(unusedFiles[0]) != "unused.golden" {
		t.Errorf("expected unused.golden, got %s", filepath.Base(unusedFiles[0]))
	}
}
