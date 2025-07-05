package gilt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoldenFile_ListFiles(t *testing.T) {
	namespace := "TestListFiles"
	
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
	testFiles := []string{"test1.golden", "test2.golden", "test3.golden"}
	for _, filename := range testFiles {
		filePath := filepath.Join(testDir, filename)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	
	// GoldenFileを作成
	goldenFile := NewGoldenFile(namespace)
	
	// ListFilesメソッドを呼び出し（まだ実装されていないのでコンパイルエラー）
	files := goldenFile.ListFiles(t, namespace)
	
	// 結果を検証
	if len(files) != 3 {
		t.Errorf("expected 3 files, got %d", len(files))
	}
	
	// ファイル名を検証
	expectedFiles := make(map[string]bool)
	for _, filename := range testFiles {
		expectedFiles[filepath.Join(testDir, filename)] = true
	}
	
	for _, file := range files {
		if !expectedFiles[file] {
			t.Errorf("unexpected file: %s", file)
		}
	}
}