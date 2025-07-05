package gilt_test

import (
	"testing"

	"github.com/kmio11/gilt"
)

// TestGoldenFile_ListFiles tests the ListFiles functionality
func TestGoldenFile_ListFiles(t *testing.T) {
	// Create a GoldenFile instance
	goldenFile := gilt.NewGoldenFile("TestGoldenFile_ListFiles")

	// Create some test data to ensure golden files exist
	golden := gilt.NewStringGolden("TestGoldenFile_ListFiles")
	golden.Assert(t, "test1", "file1", func(t *testing.T, actual, expected string) {
		if actual != expected {
			t.Errorf("Expected %q, got %q", expected, actual)
		}
	})
	golden.Assert(t, "test2", "file2", func(t *testing.T, actual, expected string) {
		if actual != expected {
			t.Errorf("Expected %q, got %q", expected, actual)
		}
	})

	// Test ListFiles
	files := goldenFile.ListFiles(t, "TestGoldenFile_ListFiles")

	// Should have at least the files we just created
	if len(files) < 2 {
		t.Errorf("Expected at least 2 files, got %d", len(files))
	}

	// Verify files contain expected names
	foundFile1, foundFile2 := false, false
	for _, file := range files {
		if file != "" {
			if file[len(file)-12:] == "file1.golden" {
				foundFile1 = true
			}
			if file[len(file)-12:] == "file2.golden" {
				foundFile2 = true
			}
		}
	}

	if !foundFile1 || !foundFile2 {
		t.Errorf("Expected to find file1.golden and file2.golden in %v", files)
	}
}
