package gilt_test

import (
	"testing"

	"github.com/kmio11/gilt"
)

func TestCommandLineFlags(t *testing.T) {
	t.Run("DefaultCleanerBehavior", func(t *testing.T) {
		// Test the default cleaner behavior
		golden := gilt.NewStringGolden("test-flags")

		// Create some test data
		golden.Assert(t, "test content", "default-test", func(t *testing.T, actual, expected string) {
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
			}
		})

		t.Log("Default cleaner behavior test completed")
	})
}

func TestFlagIntegration(t *testing.T) {
	t.Run("AutoCleanIntegration", func(t *testing.T) {
		// Create a golden instance - if clean flags are set, auto-clean should be enabled
		golden := gilt.NewStringGolden("test-flags")

		// Create test data
		testData := "test content"

		// This should work regardless of flag state
		golden.Assert(t, testData, "flag-test", func(t *testing.T, actual, expected string) {
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
			}
		})

		t.Log("Flag integration test completed successfully")
	})
}

// Example of how to set flags programmatically in tests
func TestNewCleanupInterface(t *testing.T) {
	t.Run("EnableCleanupWithOptions", func(t *testing.T) {
		// Test EnableCleanup with custom options
		golden := gilt.NewStringGolden("test-custom")
		gilt.SetupCleanup(t, golden, gilt.WithDryRun(), gilt.WithBackup())

		// Test with cleanup enabled
		testData := "custom cleanup test content"
		golden.Assert(t, testData, "custom-cleanup-test", func(t *testing.T, actual, expected string) {
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
			}
		})

		t.Log("Custom cleanup test completed")
	})

	t.Run("EnableCleanupSimple", func(t *testing.T) {
		// Test simple EnableCleanup
		golden := gilt.NewStringGolden("test-simple")
		gilt.SetupCleanup(t, golden)

		testData := "simple cleanup test content"
		golden.Assert(t, testData, "simple-cleanup-test", func(t *testing.T, actual, expected string) {
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
			}
		})

		t.Log("Simple cleanup test completed")
	})
}
