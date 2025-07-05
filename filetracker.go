package gilt

import (
	"sync"
	"testing"
)

// FileTracker tracks which files have been accessed during testing.
type FileTracker struct {
	usedFiles map[string]bool
	mutex     sync.RWMutex
}

// NewFileTracker creates a new FileTracker instance.
func NewFileTracker() *FileTracker {
	return &FileTracker{
		usedFiles: make(map[string]bool),
	}
}

// TrackUsage marks a file as used by recording its path.
func (ft *FileTracker) TrackUsage(t *testing.T, filePath string) {
	t.Helper()
	ft.mutex.Lock()
	defer ft.mutex.Unlock()
	ft.usedFiles[filePath] = true
}

// IsUsed checks whether a file has been marked as used.
func (ft *FileTracker) IsUsed(filePath string) bool {
	ft.mutex.RLock()
	defer ft.mutex.RUnlock()
	return ft.usedFiles[filePath]
}

// GetUnusedFiles returns a list of files from the provided list that have not been marked as used.
func (ft *FileTracker) GetUnusedFiles(allFiles []string) []string {
	ft.mutex.RLock()
	defer ft.mutex.RUnlock()

	var unused []string
	for _, file := range allFiles {
		if !ft.usedFiles[file] {
			unused = append(unused, file)
		}
	}
	return unused
}

// Reset clears all tracked file usage records.
func (ft *FileTracker) Reset() {
	ft.mutex.Lock()
	defer ft.mutex.Unlock()
	ft.usedFiles = make(map[string]bool)
}
