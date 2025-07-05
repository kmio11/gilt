package gilt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type CleanOptions struct {
	DryRun      bool
	Backup      bool
	Interactive bool
	Patterns    []string
}

type CleanResult struct {
	DeletedFiles     []string
	WouldDeleteFiles []string
	BackupDir        string
	LogMessages      []string
	Error            error
}

func CleanUnusedGoldenFiles(goldenFile *GoldenFile, t *testing.T, namespace string, options CleanOptions) CleanResult {
	t.Helper()
	
	result := CleanResult{
		DeletedFiles:     []string{},
		WouldDeleteFiles: []string{},
		LogMessages:      []string{},
	}
	
	// Listerインターフェースを確認
	lister, ok := interface{}(goldenFile).(Lister)
	if !ok {
		result.Error = fmt.Errorf("goldenFile does not implement Lister interface")
		return result
	}
	
	// 全ファイルを取得
	allFiles := lister.ListFiles(t, namespace)
	
	// 未使用ファイルを取得
	unusedFiles := GetUnusedFiles(allFiles)
	
	// パターンフィルタリングを適用
	filteredFiles := applyPatternFiltering(unusedFiles, options.Patterns)
	
	// ログに記録
	result.LogMessages = append(result.LogMessages, fmt.Sprintf("Found %d total files", len(allFiles)))
	result.LogMessages = append(result.LogMessages, fmt.Sprintf("Found %d unused files", len(unusedFiles)))
	result.LogMessages = append(result.LogMessages, fmt.Sprintf("Found %d files to process after filtering", len(filteredFiles)))
	
	// Dry-runの場合
	if options.DryRun {
		result.WouldDeleteFiles = filteredFiles
		for _, file := range filteredFiles {
			result.LogMessages = append(result.LogMessages, fmt.Sprintf("Would delete: %s", file))
		}
		return result
	}
	
	// バックアップディレクトリを作成
	var backupDir string
	if options.Backup {
		backupDir = filepath.Join("testdata", namespace, ".backup")
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			result.Error = fmt.Errorf("failed to create backup directory: %v", err)
			return result
		}
		result.BackupDir = backupDir
	}
	
	// 未使用ファイルを削除
	for _, file := range filteredFiles {
		// バックアップを作成
		if options.Backup {
			backupFile := filepath.Join(backupDir, filepath.Base(file))
			if err := copyFile(file, backupFile); err != nil {
				result.Error = fmt.Errorf("failed to backup file %s: %v", file, err)
				return result
			}
			result.LogMessages = append(result.LogMessages, fmt.Sprintf("Backed up: %s -> %s", file, backupFile))
		}
		
		// ファイルを削除
		if err := os.Remove(file); err != nil {
			result.Error = fmt.Errorf("failed to remove file %s: %v", file, err)
			return result
		}
		
		result.DeletedFiles = append(result.DeletedFiles, file)
		result.LogMessages = append(result.LogMessages, fmt.Sprintf("Deleted: %s", file))
	}
	
	return result
}

func applyPatternFiltering(files []string, patterns []string) []string {
	if len(patterns) == 0 {
		return files
	}
	
	var filtered []string
	for _, file := range files {
		basename := filepath.Base(file)
		protected := false
		
		for _, pattern := range patterns {
			if matched, _ := filepath.Match(pattern, basename); matched {
				protected = true
				break
			}
		}
		
		if !protected {
			filtered = append(filtered, file)
		}
	}
	
	return filtered
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	
	return os.WriteFile(dst, data, 0644)
}