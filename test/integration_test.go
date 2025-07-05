package gilt_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kmio11/gilt"
)

func TestEndToEnd_CleanUnusedFiles(t *testing.T) {
	namespace := "TestEndToEnd_CleanUnusedFiles"

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
	testFiles := []string{
		"active_test1.golden",
		"active_test2.golden",
		"old_unused1.golden",
		"old_unused2.golden",
		"important_config.golden",
	}
	for _, filename := range testFiles {
		filePath := filepath.Join(testDir, filename)
		err := os.WriteFile(filePath, []byte("test content for "+filename), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// GoldenFileを直接使用してファイルを読み込み（使用状況を記録）
	goldenFile := gilt.NewGoldenFile(namespace)

	// active_test1とactive_test2のみ読み込み（使用状況を記録）
	data1 := goldenFile.Read(t, "active_test1")
	if string(data1) != "test content for active_test1.golden" {
		t.Errorf("expected correct content for active_test1, got %s", string(data1))
	}

	data2 := goldenFile.Read(t, "active_test2")
	if string(data2) != "test content for active_test2.golden" {
		t.Errorf("expected correct content for active_test2, got %s", string(data2))
	}

	// 未使用ファイルを削除（important_configは保護）
	result := gilt.CleanUnusedGoldenFiles(goldenFile, t, namespace, gilt.CleanOptions{
		DryRun:   false,
		Backup:   true,
		Patterns: []string{"important_*"}, // important_configを保護
	})

	// 結果を検証
	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}

	// 使用されたファイルが残っていることを確認
	activeFiles := []string{"active_test1.golden", "active_test2.golden"}
	for _, filename := range activeFiles {
		filePath := filepath.Join(testDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("%s should still exist", filename)
		}
	}

	// 保護されたファイルが残っていることを確認
	importantFile := filepath.Join(testDir, "important_config.golden")
	if _, err := os.Stat(importantFile); os.IsNotExist(err) {
		t.Error("important_config.golden should be protected")
	}

	// 未使用ファイルが削除されていることを確認
	unusedFiles := []string{"old_unused1.golden", "old_unused2.golden"}
	for _, filename := range unusedFiles {
		filePath := filepath.Join(testDir, filename)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Errorf("%s should be deleted", filename)
		}
	}

	// バックアップが作成されていることを確認
	backupDir := filepath.Join("testdata", namespace, ".backup")
	for _, filename := range unusedFiles {
		backupFile := filepath.Join(backupDir, filename)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			t.Errorf("backup for %s should exist", filename)
		}
	}

	// 削除されたファイル数を確認
	if len(result.DeletedFiles) != 2 {
		t.Errorf("expected 2 deleted files, got %d", len(result.DeletedFiles))
	}
}

func TestEndToEnd_WithExistingGoldenFiles(t *testing.T) {
	namespace := "TestEndToEnd_WithExistingGoldenFiles"

	// まずゴールデンファイルを作成するために手動でファイルを作成
	testDir := filepath.Join("testdata", namespace, "golden")
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// 使用されるファイルを作成
	usedFiles := map[string]string{
		"hello_world.golden":  "Hello, world!",
		"hello_gopher.golden": "Hello, gopher!",
	}

	for filename, content := range usedFiles {
		filePath := filepath.Join(testDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// GoldenFileを使用してファイルを読み込み（使用状況を記録）
	goldenFile := gilt.NewGoldenFile(namespace)

	for name := range usedFiles {
		basename := name[:len(name)-len(".golden")]
		data := goldenFile.Read(t, basename)
		if len(data) == 0 {
			t.Errorf("failed to read %s", name)
		}
	}

	// 手動で未使用ファイルを作成
	unusedFile := filepath.Join(testDir, "unused_test.golden")
	err = os.WriteFile(unusedFile, []byte("unused content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// クリーンアップ
	t.Cleanup(func() {
		os.RemoveAll(filepath.Join("testdata", namespace))
	})

	// 未使用ファイルをdry-runで確認
	result := gilt.CleanUnusedGoldenFiles(goldenFile, t, namespace, gilt.CleanOptions{
		DryRun: true,
	})

	// 結果を検証
	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}

	// dry-runでは1つの未使用ファイルが検出されるはず
	if len(result.WouldDeleteFiles) != 1 {
		t.Errorf("expected 1 file would be deleted, got %d", len(result.WouldDeleteFiles))
	}

	// ファイルはまだ存在するはず
	if _, err := os.Stat(unusedFile); os.IsNotExist(err) {
		t.Error("unused file should still exist after dry-run")
	}
}
