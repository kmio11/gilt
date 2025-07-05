package gilt

import (
	"sync"
	"testing"
)

var (
	globalFileTracker *FileTracker
	trackingEnabled   bool
	trackingMutex     sync.RWMutex
)

type FileTracker struct {
	usedFiles map[string]bool
	mutex     sync.RWMutex
}

func NewFileTracker() *FileTracker {
	return &FileTracker{
		usedFiles: make(map[string]bool),
	}
}

func (ft *FileTracker) TrackUsage(t *testing.T, filePath string) {
	t.Helper()
	ft.mutex.Lock()
	defer ft.mutex.Unlock()
	ft.usedFiles[filePath] = true
}

func (ft *FileTracker) IsUsed(filePath string) bool {
	ft.mutex.RLock()
	defer ft.mutex.RUnlock()
	return ft.usedFiles[filePath]
}

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

func (ft *FileTracker) Reset() {
	ft.mutex.Lock()
	defer ft.mutex.Unlock()
	ft.usedFiles = make(map[string]bool)
}

func EnableFileTracking() {
	trackingMutex.Lock()
	defer trackingMutex.Unlock()
	trackingEnabled = true
	if globalFileTracker == nil {
		globalFileTracker = NewFileTracker()
	}
}

func DisableFileTracking() {
	trackingMutex.Lock()
	defer trackingMutex.Unlock()
	trackingEnabled = false
}

func ResetGlobalFileTracker() {
	trackingMutex.Lock()
	defer trackingMutex.Unlock()
	if globalFileTracker != nil {
		globalFileTracker.Reset()
	}
}

func TrackFileUsage(t *testing.T, filePath string) {
	trackingMutex.RLock()
	defer trackingMutex.RUnlock()
	if trackingEnabled && globalFileTracker != nil {
		globalFileTracker.TrackUsage(t, filePath)
	}
}

func GetUnusedFiles(allFiles []string) []string {
	trackingMutex.RLock()
	defer trackingMutex.RUnlock()
	if trackingEnabled && globalFileTracker != nil {
		return globalFileTracker.GetUnusedFiles(allFiles)
	}
	return []string{}
}