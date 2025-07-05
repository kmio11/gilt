package gilt

import (
	"os"
	"path/filepath"
	"testing"
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
	defer func() {
		os.RemoveAll(filepath.Join("testdata", namespace))
	}()
	
	// テスト用のファイルを作成
	testFiles := []string{"used.golden", "unused.golden"}
	for _, filename := range testFiles {
		filePath := filepath.Join(testDir, filename)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	
	// GlobalFileTrackerを使用してファイル追跡を有効化
	ResetGlobalFileTracker()
	EnableFileTracking()
	
	// GoldenFileを作成
	goldenFile := NewGoldenFile(namespace)
	
	// ファイルを読み込み（自動的に追跡される）
	data := goldenFile.Read(t, "used")
	if string(data) != "test content" {
		t.Errorf("expected 'test content', got '%s'", string(data))
	}
	
	// ファイル一覧を取得
	allFiles := goldenFile.ListFiles(t, namespace)
	
	// 未使用ファイルを取得
	unusedFiles := GetUnusedFiles(allFiles)
	
	// 検証
	if len(unusedFiles) != 1 {
		t.Errorf("expected 1 unused file, got %d", len(unusedFiles))
	}
	if len(unusedFiles) > 0 && filepath.Base(unusedFiles[0]) != "unused.golden" {
		t.Errorf("expected unused.golden, got %s", filepath.Base(unusedFiles[0]))
	}
	
	// ファイル追跡を無効化
	DisableFileTracking()
}