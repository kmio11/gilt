package gilt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanUnusedGoldenFiles(t *testing.T) {
	namespace := "TestCleanUnusedGoldenFiles"
	
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
	testFiles := []string{"used.golden", "unused1.golden", "unused2.golden"}
	for _, filename := range testFiles {
		filePath := filepath.Join(testDir, filename)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	
	// ファイル追跡を有効化
	ResetGlobalFileTracker()
	EnableFileTracking()
	defer DisableFileTracking()
	
	// GoldenFileを作成
	goldenFile := NewGoldenFile(namespace)
	
	// 1つのファイルのみ使用
	goldenFile.Read(t, "used")
	
	// 未使用ファイルを削除（まだ実装されていない）
	options := CleanOptions{
		DryRun: false,
		Backup: false,
	}
	
	result := CleanUnusedGoldenFiles(goldenFile, t, namespace, options)
	
	// 結果を検証
	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}
	
	if len(result.DeletedFiles) != 2 {
		t.Errorf("expected 2 deleted files, got %d", len(result.DeletedFiles))
	}
	
	// 使用されたファイルが残っていることを確認
	usedFilePath := filepath.Join(testDir, "used.golden")
	if _, err := os.Stat(usedFilePath); os.IsNotExist(err) {
		t.Error("used.golden should still exist")
	}
	
	// 未使用ファイルが削除されていることを確認
	unused1FilePath := filepath.Join(testDir, "unused1.golden")
	if _, err := os.Stat(unused1FilePath); !os.IsNotExist(err) {
		t.Error("unused1.golden should be deleted")
	}
	
	unused2FilePath := filepath.Join(testDir, "unused2.golden")
	if _, err := os.Stat(unused2FilePath); !os.IsNotExist(err) {
		t.Error("unused2.golden should be deleted")
	}
}

func TestCleanUnusedGoldenFiles_DryRun(t *testing.T) {
	namespace := "TestCleanUnusedGoldenFiles_DryRun"
	
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
	
	// ファイル追跡を有効化
	ResetGlobalFileTracker()
	EnableFileTracking()
	defer DisableFileTracking()
	
	// GoldenFileを作成
	goldenFile := NewGoldenFile(namespace)
	
	// 1つのファイルのみ使用
	goldenFile.Read(t, "used")
	
	// Dry-runで未使用ファイルを確認
	options := CleanOptions{
		DryRun: true,
		Backup: false,
	}
	
	result := CleanUnusedGoldenFiles(goldenFile, t, namespace, options)
	
	// 結果を検証
	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}
	
	if len(result.WouldDeleteFiles) != 1 {
		t.Errorf("expected 1 file would be deleted, got %d", len(result.WouldDeleteFiles))
	}
	
	// ファイルが削除されていないことを確認
	unusedFilePath := filepath.Join(testDir, "unused.golden")
	if _, err := os.Stat(unusedFilePath); os.IsNotExist(err) {
		t.Error("unused.golden should still exist in dry-run mode")
	}
}

func TestCleanUnusedGoldenFiles_Backup(t *testing.T) {
	namespace := "TestCleanUnusedGoldenFiles_Backup"
	
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
	
	// ファイル追跡を有効化
	ResetGlobalFileTracker()
	EnableFileTracking()
	defer DisableFileTracking()
	
	// GoldenFileを作成
	goldenFile := NewGoldenFile(namespace)
	
	// 1つのファイルのみ使用
	goldenFile.Read(t, "used")
	
	// バックアップ付きで未使用ファイルを削除
	options := CleanOptions{
		DryRun: false,
		Backup: true,
	}
	
	result := CleanUnusedGoldenFiles(goldenFile, t, namespace, options)
	
	// 結果を検証
	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}
	
	if len(result.DeletedFiles) != 1 {
		t.Errorf("expected 1 deleted file, got %d", len(result.DeletedFiles))
	}
	
	// バックアップファイルが作成されていることを確認
	backupDir := filepath.Join("testdata", namespace, ".backup")
	backupFile := filepath.Join(backupDir, "unused.golden")
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Error("backup file should exist")
	}
}