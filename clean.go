package gilt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// CleanOptions configures the behavior of the file cleaning operation.
type CleanOptions struct {
	DryRun      bool     // If true, only simulate the cleaning without actually deleting files
	Backup      bool     // If true, create backups of deleted files
	Interactive bool     // If true, prompt for confirmation before deleting files
	Patterns    []string // File patterns to exclude from deletion
}

// CleanResult contains the results of a file cleaning operation.
type CleanResult struct {
	DeletedFiles     []string // Files that were actually deleted
	WouldDeleteFiles []string // Files that would be deleted in a dry run
	BackupDir        string   // Directory where backups were created
	LogMessages      []string // Log messages from the cleaning operation
	Error            error    // Any error that occurred during cleaning
}

// CleanUnusedGoldenFiles removes unused golden files based on the provided options.
// It identifies files that haven't been accessed during the test and optionally creates backups before deletion.
func CleanUnusedGoldenFiles(goldenFile *GoldenFile, t *testing.T, namespace string, options CleanOptions) CleanResult {
	t.Helper()

	result := CleanResult{
		DeletedFiles:     []string{},
		WouldDeleteFiles: []string{},
		LogMessages:      []string{},
	}

	// Check if goldenFile implements Lister interface
	lister, ok := interface{}(goldenFile).(Lister)
	if !ok {
		result.Error = fmt.Errorf("goldenFile does not implement Lister interface")
		return result
	}

	// Get all files
	allFiles := lister.ListFiles(t, namespace)

	// Get unused files
	unusedFiles := goldenFile.GetUnusedFiles(t, namespace)

	// Apply pattern filtering
	filteredFiles := applyPatternFiltering(unusedFiles, options.Patterns)

	// Log results
	result.LogMessages = append(result.LogMessages, fmt.Sprintf("Found %d total files", len(allFiles)))
	result.LogMessages = append(result.LogMessages, fmt.Sprintf("Found %d unused files", len(unusedFiles)))
	result.LogMessages = append(result.LogMessages, fmt.Sprintf("Found %d files to process after filtering", len(filteredFiles)))

	// Handle dry run mode
	if options.DryRun {
		result.WouldDeleteFiles = filteredFiles
		for _, file := range filteredFiles {
			result.LogMessages = append(result.LogMessages, fmt.Sprintf("Would delete: %s", file))
		}
		return result
	}

	// Create backup directory
	var backupDir string
	if options.Backup {
		backupDir = filepath.Join("testdata", namespace, ".backup")
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			result.Error = fmt.Errorf("failed to create backup directory: %v", err)
			return result
		}
		result.BackupDir = backupDir
	}

	// Delete unused files
	for _, file := range filteredFiles {
		// Create backup
		if options.Backup {
			backupFile := filepath.Join(backupDir, filepath.Base(file))
			if err := copyFile(file, backupFile); err != nil {
				result.Error = fmt.Errorf("failed to backup file %s: %v", file, err)
				return result
			}
			result.LogMessages = append(result.LogMessages, fmt.Sprintf("Backed up: %s -> %s", file, backupFile))
		}

		// Delete file
		if err := os.Remove(file); err != nil {
			result.Error = fmt.Errorf("failed to remove file %s: %v", file, err)
			return result
		}

		result.DeletedFiles = append(result.DeletedFiles, file)
		result.LogMessages = append(result.LogMessages, fmt.Sprintf("Deleted: %s", file))
	}

	return result
}

// applyPatternFiltering filters out files that match any of the provided patterns.
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

// copyFile copies a file from source to destination.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}
