package gilt_test

import (
	"testing"

	"github.com/kmio11/gilt"
)

// TestSaferCleanupPattern demonstrates the improved cleanup approach
func TestSaferCleanupPattern(t *testing.T) {
	golden := gilt.NewStringGolden("safer-test")

	// Enable cleanup at the test function level
	gilt.SetupCleanup(t, golden, gilt.WithDryRun())

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"safe1", "test1", "test1"},
		{"safe2", "test2", "test2"},
		{"safe3", "test3", "test3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Now Assert does not trigger cleanup - it's handled at the test function level
			golden.Assert(t, tc.input, tc.name, func(t *testing.T, actual, expected string) {
				if actual != expected {
					t.Errorf("Expected %q, got %q", expected, actual)
				}
			})
		})
	}
	// Cleanup will be called once at the end of this test function
}

// TestSaferParallelPattern demonstrates safer parallel test cleanup
func TestSaferParallelPattern(t *testing.T) {
	golden := gilt.NewStringGolden("safer-parallel")

	// Enable cleanup at the test function level
	gilt.SetupCleanup(t, golden, gilt.WithDryRun())

	testCases := []struct {
		name  string
		input string
	}{
		{"safe_parallel1", "test1"},
		{"safe_parallel2", "test2"},
		{"safe_parallel3", "test3"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			golden.Assert(t, tc.input, tc.name, func(t *testing.T, actual, expected string) {
				if actual != expected {
					t.Errorf("Expected %q, got %q", expected, actual)
				}
			})
		})
	}
	// Cleanup will be called once after all parallel subtests complete
}

// TestWithoutCleanup demonstrates tests without any cleanup
func TestWithoutCleanup(t *testing.T) {
	golden := gilt.NewStringGolden("no-cleanup")

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"no_cleanup1", "test1", "test1"},
		{"no_cleanup2", "test2", "test2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			golden.Assert(t, tc.input, tc.name, func(t *testing.T, actual, expected string) {
				if actual != expected {
					t.Errorf("Expected %q, got %q", expected, actual)
				}
			})
		})
	}
	// No cleanup - files remain
}

// TestCommandLineFlagIntegration demonstrates how command line flags work with the safer approach
func TestCommandLineFlagIntegration(t *testing.T) {
	// This golden instance will use command line flags for cleanup decisions
	golden := gilt.NewStringGolden("flag-integration")

	// Enable cleanup - this will respect command line flags
	gilt.SetupCleanup(t, golden)

	testCases := []string{"flag1", "flag2", "flag3"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			golden.Assert(t, tc, tc, func(t *testing.T, actual, expected string) {
				if actual != expected {
					t.Errorf("Expected %q, got %q", expected, actual)
				}
			})
		})
	}
	// Cleanup behavior depends on command line flags:
	// go test -clean          -> cleanup enabled
	// go test -clean-dry-run  -> dry run cleanup
	// go test -clean-backup   -> cleanup with backup
	// go test                 -> no cleanup
}
