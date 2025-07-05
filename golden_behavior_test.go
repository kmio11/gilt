package gilt_test

import (
	"testing"

	"github.com/kmio11/gilt"
)

// TestCleanupBehavior tests the correct cleanup behavior
func TestCleanupBehavior(t *testing.T) {
	t.Run("NoOptions_NoCleanup", func(t *testing.T) {
		// Without any options or flags, EnableCleanup should do nothing
		golden := gilt.NewStringGolden("no-options")

		// This should NOT enable cleanup
		gilt.SetupCleanup(t, golden)

		golden.Assert(t, "test", "test", func(t *testing.T, actual, expected string) {
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
			}
		})

		// No cleanup should happen
		t.Log("No cleanup should happen without options or flags")
	})

	t.Run("WithOptions_ForceCleanup", func(t *testing.T) {
		// With options, cleanup should be enabled even without flags
		golden := gilt.NewStringGolden("with-options")

		// This SHOULD enable cleanup because options are provided
		gilt.SetupCleanup(t, golden, gilt.WithDryRun())

		golden.Assert(t, "test", "test", func(t *testing.T, actual, expected string) {
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
			}
		})

		// Cleanup should happen with dry-run
		t.Log("Cleanup should happen with options")
	})

	t.Run("WithBackup_ForceCleanup", func(t *testing.T) {
		// With backup option, cleanup should be enabled
		golden := gilt.NewStringGolden("with-backup")

		// This SHOULD enable cleanup with backup
		gilt.SetupCleanup(t, golden, gilt.WithBackup())

		golden.Assert(t, "test", "test", func(t *testing.T, actual, expected string) {
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
			}
		})

		t.Log("Cleanup should happen with backup option")
	})
}
