package gilt_test

import (
	"testing"

	"github.com/kmio11/gilt"
)

// TestSimplifiedAPI demonstrates the new simplified API
func TestSimplifiedAPI(t *testing.T) {
	t.Run("DefaultBehavior", func(t *testing.T) {
		// Default behavior - no explicit cleanup configuration
		// Cleanup is controlled by command line flags only
		golden := gilt.NewStringGolden("simplified-default")

		golden.Assert(t, "default behavior", "test1", func(t *testing.T, actual, expected string) {
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
			}
		})

		// go test -clean will enable cleanup
		// go test will not cleanup
		t.Log("Default behavior test completed")
	})

	t.Run("ExplicitCleanup", func(t *testing.T) {
		// Explicit cleanup - overrides command line flags
		golden := gilt.NewStringGolden("simplified-explicit")
		gilt.SetupCleanup(t, golden) // Force enable cleanup

		golden.Assert(t, "explicit cleanup", "test2", func(t *testing.T, actual, expected string) {
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
			}
		})

		t.Log("Explicit cleanup test completed")
	})

	t.Run("CleanupWithOptions", func(t *testing.T) {
		// Cleanup with custom options
		golden := gilt.NewStringGolden("simplified-options")
		gilt.SetupCleanup(t, golden,
			gilt.WithDryRun(),
			gilt.WithBackup(),
			gilt.WithPatterns("*.keep", "important_*"),
		)

		golden.Assert(t, "cleanup with options", "test3", func(t *testing.T, actual, expected string) {
			if actual != expected {
				t.Errorf("Expected %q, got %q", expected, actual)
			}
		})

		t.Log("Cleanup with options test completed")
	})

	t.Run("TableDrivenSafe", func(t *testing.T) {
		// Table-driven test with safe cleanup timing
		golden := gilt.NewStringGolden("simplified-table")
		gilt.SetupCleanup(t, golden, gilt.WithDryRun())

		testCases := []struct {
			name  string
			input string
		}{
			{"case1", "test1"},
			{"case2", "test2"},
			{"case3", "test3"},
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

		// Cleanup happens once at the end of this test function
		t.Log("Table-driven safe test completed")
	})
}
