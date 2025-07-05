package gilt_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kmio11/gilt"
)

func TestCleanUnusedGoldenFiles_WithPatterns(t *testing.T) {
	namespace := "TestCleanUnusedGoldenFiles_WithPatterns"

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
	testFiles := []string{"used.golden", "unused.golden", "important.golden", "temp.golden"}
	for _, filename := range testFiles {
		filePath := filepath.Join(testDir, filename)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// GoldenFileを作成
	goldenFile := gilt.NewGoldenFile(namespace)

	// 1つのファイルのみ使用
	goldenFile.Read(t, "used")

	// 保護パターンを指定して削除
	options := gilt.CleanOptions{
		DryRun:   false,
		Backup:   false,
		Patterns: []string{"important*"}, // importantで始まるファイルは保護
	}

	result := gilt.CleanUnusedGoldenFiles(goldenFile, t, namespace, options)

	// 結果を検証
	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}

	// important.goldenは保護されているので削除されないはず
	importantFilePath := filepath.Join(testDir, "important.golden")
	if _, err := os.Stat(importantFilePath); os.IsNotExist(err) {
		t.Error("important.golden should be protected by pattern")
	}

	// temp.goldenは削除されているはず
	tempFilePath := filepath.Join(testDir, "temp.golden")
	if _, err := os.Stat(tempFilePath); !os.IsNotExist(err) {
		t.Error("temp.golden should be deleted")
	}
}

func TestCleanUnusedGoldenFiles_WithLogging(t *testing.T) {
	namespace := "TestCleanUnusedGoldenFiles_WithLogging"

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

	// GoldenFileを作成
	goldenFile := gilt.NewGoldenFile(namespace)

	// 1つのファイルのみ使用
	goldenFile.Read(t, "used")

	// ログ出力を有効化して削除
	options := gilt.CleanOptions{
		DryRun: false,
		Backup: false,
	}

	result := gilt.CleanUnusedGoldenFiles(goldenFile, t, namespace, options)

	// 結果を検証
	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}

	// ログメッセージが含まれていることを確認
	if len(result.LogMessages) == 0 {
		t.Error("expected log messages to be recorded")
	}

	// 削除されたファイルがログに含まれていることを確認
	found := false
	for _, msg := range result.LogMessages {
		if strings.Contains(msg, "unused.golden") && strings.Contains(msg, "Deleted") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected deletion log message for unused.golden, got: %v", result.LogMessages)
	}
}
