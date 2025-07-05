package gilt_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kmio11/gilt"
)

func TestGolden_CleanUnusedFiles(t *testing.T) {
	namespace := "TestGolden_CleanUnusedFiles"

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

	// Goldenインスタンスを作成
	golden := gilt.NewStringGolden(namespace)

	// 1つのファイルのみ使用
	golden.Assert(t, "test content", "used", func(t *testing.T, actual, expected string) {
		if actual != expected {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})

	// Goldenインスタンスから直接クリーン実行
	result := golden.CleanUnusedFiles(t, gilt.CleanOptions{
		DryRun: false,
		Backup: true,
	})

	// 結果を検証
	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}

	if len(result.DeletedFiles) != 1 {
		t.Errorf("expected 1 deleted file, got %d", len(result.DeletedFiles))
	}
}

func TestWithAutoClean(t *testing.T) {
	namespace := "TestWithAutoClean"

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

	// 自動クリーン機能付きGoldenインスタンスを作成
	golden := gilt.NewStringGolden(namespace, gilt.WithAutoClean[string, string](gilt.CleanOptions{
		DryRun: false,
		Backup: true,
	}))

	// 1つのファイルのみ使用
	golden.Assert(t, "test content", "used", func(t *testing.T, actual, expected string) {
		if actual != expected {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	})

	// テスト終了時に自動でクリーンが実行される（t.Cleanup使用）
	// t.Cleanupは手動で実行
	// NOTE: 実際のテストではt.Cleanupが自動実行されるが、
	// テスト内で確認するため手動で実行

	// 手動でクリーンアップ関数を実行してテスト
	if goldenFile, ok := golden.GetGoldenFile().(*gilt.GoldenFile); ok {
		result := gilt.CleanUnusedGoldenFiles(goldenFile, t, namespace, gilt.CleanOptions{
			DryRun: false,
			Backup: true,
		})
		if result.Error != nil {
			t.Errorf("cleanup failed: %v", result.Error)
		}
		if len(result.DeletedFiles) == 0 {
			t.Error("expected some files to be deleted")
		}
	}
}
