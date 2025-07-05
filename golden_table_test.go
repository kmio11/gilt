package gilt_test

import (
	"testing"

	"github.com/kmio11/gilt"
)

// TestTableDrivenPattern demonstrates the cleanup timing issue
func TestTableDrivenPattern(t *testing.T) {
	// This test demonstrates the safer cleanup implementation
	golden := gilt.NewStringGolden("table-test")

	// Enable cleanup at test function level to avoid timing issues
	gilt.SetupCleanup(t, golden, gilt.WithDryRun())

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"case1", "test1", "test1"},
		{"case2", "test2", "test2"},
		{"case3", "test3", "test3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Each call to Assert will trigger cleanup registration
			// This demonstrates the timing issue
			golden.Assert(t, tc.input, tc.name, func(t *testing.T, actual, expected string) {
				if actual != expected {
					t.Errorf("Expected %q, got %q", expected, actual)
				}
			})
		})
	}
}

// TestParallelPattern demonstrates potential issues with parallel tests
func TestParallelPattern(t *testing.T) {
	golden := gilt.NewStringGolden("parallel-test")
	gilt.SetupCleanup(t, golden, gilt.WithDryRun())

	testCases := []struct {
		name  string
		input string
	}{
		{"parallel1", "test1"},
		{"parallel2", "test2"},
		{"parallel3", "test3"},
	}

	for _, tc := range testCases {
		tc := tc // capture loop variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			golden.Assert(t, tc.input, tc.name, func(t *testing.T, actual, expected string) {
				if actual != expected {
					t.Errorf("Expected %q, got %q", expected, actual)
				}
			})
		})
	}
}

// TestAutoCleanupPattern demonstrates the recommended approach
func TestAutoCleanupPattern(t *testing.T) {
	golden := gilt.NewStringGolden("auto-test")

	// Enable cleanup at the test function level
	gilt.SetupCleanup(t, golden, gilt.WithDryRun())

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"auto1", "test1", "test1"},
		{"auto2", "test2", "test2"},
		{"auto3", "test3", "test3"},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			golden.Assert(t, tc.input, tc.name, func(t *testing.T, actual, expected string) {
				if actual != expected {
					t.Errorf("Expected %q, got %q", expected, actual)
				}
			})
		})
	}
	// Cleanup happens automatically at the end
}
