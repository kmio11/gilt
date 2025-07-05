package gilt_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kmio11/gilt"
)

func TestFileTracker(t *testing.T) {
	namespace := "TestFileTracker"

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
	testFiles := []string{"used1.golden", "used2.golden", "unused.golden"}
	for _, filename := range testFiles {
		filePath := filepath.Join(testDir, filename)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// FileTrackerを作成（まだ実装されていない）
	tracker := gilt.NewFileTracker()

	// ファイル使用を記録
	tracker.TrackUsage(t, filepath.Join(testDir, "used1.golden"))
	tracker.TrackUsage(t, filepath.Join(testDir, "used2.golden"))

	// 使用状況を確認
	if !tracker.IsUsed(filepath.Join(testDir, "used1.golden")) {
		t.Error("used1.golden should be marked as used")
	}
	if !tracker.IsUsed(filepath.Join(testDir, "used2.golden")) {
		t.Error("used2.golden should be marked as used")
	}
	if tracker.IsUsed(filepath.Join(testDir, "unused.golden")) {
		t.Error("unused.golden should not be marked as used")
	}

	// 未使用ファイルを取得
	allFiles := []string{
		filepath.Join(testDir, "used1.golden"),
		filepath.Join(testDir, "used2.golden"),
		filepath.Join(testDir, "unused.golden"),
	}

	unusedFiles := tracker.GetUnusedFiles(allFiles)
	if len(unusedFiles) != 1 {
		t.Errorf("expected 1 unused file, got %d", len(unusedFiles))
	}
	if len(unusedFiles) > 0 && unusedFiles[0] != filepath.Join(testDir, "unused.golden") {
		t.Errorf("expected unused.golden, got %s", unusedFiles[0])
	}
}

func TestFileTracker_Reset(t *testing.T) {
	tracker := gilt.NewFileTracker()

	// ファイル使用を記録
	tracker.TrackUsage(t, "test.golden")

	// 使用状況を確認
	if !tracker.IsUsed("test.golden") {
		t.Error("test.golden should be marked as used")
	}

	// リセット
	tracker.Reset()

	// リセット後は使用されていないはず
	if tracker.IsUsed("test.golden") {
		t.Error("test.golden should not be marked as used after reset")
	}
}
