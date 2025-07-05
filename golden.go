package gilt

import (
	"os"
	"testing"
)

type (
	// ReadWriter combines Reader and Writer interfaces for golden file operations.
	ReadWriter interface {
		Reader
		Writer
	}
	// Reader provides methods to read golden file data.
	Reader interface {
		Read(t *testing.T, name string) []byte
	}
	// Writer provides methods to write golden file data.
	Writer interface {
		Write(t *testing.T, name string, data []byte)
	}
	// Opener provides methods to open golden files.
	Opener interface {
		Open(t *testing.T, name string) *os.File
	}
	// Lister provides methods to list files in a namespace.
	Lister interface {
		ListFiles(t *testing.T, namespace string) []string
	}
)

// IsUpdater determines whether golden files should be updated.
type IsUpdater interface {
	IsUpdate(t *testing.T, name string) bool
}

// Saver converts actual values to bytes for storage in golden files.
type Saver[TActual any] interface {
	Save(t *testing.T, actual TActual, name string, writer Writer)
}

// Loader reads and converts golden file data to expected values.
type Loader[TExpected any] interface {
	Load(t *testing.T, name string, reader Reader) TExpected
}

// Golden orchestrates golden file testing with generic types for actual and expected values.
type Golden[TActual, TExpected any] struct {
	updater    IsUpdater
	saver      Saver[TActual]
	loader     Loader[TExpected]
	goldenFile ReadWriter
	namespace  string
}

// Option configures a Golden instance with functional options pattern.
type Option[TActual, TExpected any] func(*Golden[TActual, TExpected])

// New creates a new Golden instance with the specified namespace and options.
// It uses JSON serialization by default for both saving and loading.
func New[TActual any, TExpected any](
	namespace string,
	opts ...Option[TActual, TExpected],
) *Golden[TActual, TExpected] {
	g := &Golden[TActual, TExpected]{
		goldenFile: NewGoldenFile(namespace),
		updater:    IsUpdaterFunc(isUpdateFunc),
		saver:      SaverFunc[TActual](SaveAsJSON[TActual]),
		loader:     LoaderFunc[TExpected](LoadAndUnmarshalJSON[TExpected]),
		namespace:  namespace,
	}
	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Assert compares actual values against golden files, updating them if needed.
// If the updater indicates an update is needed, it saves the actual value.
// Otherwise, it loads the expected value and runs the assertion function.
func (g *Golden[TActual, TExpected]) Assert(
	t *testing.T, actual TActual, name string,
	assertion func(t *testing.T, actual TActual, expected TExpected),
) {
	t.Helper()

	// Note: Auto-cleanup is NOT performed in Assert to avoid issues with table-driven tests.
	// Use gilt.SetupCleanup(t, golden) at the test function level for cleanup configuration.

	if g.updater.IsUpdate(t, name) {
		g.saver.Save(t, actual, name, g.goldenFile)
		return
	}

	expectedData := g.loader.Load(t, name, g.goldenFile)
	assertion(t, actual, expectedData)
}

// NewBytesGolden creates a Golden instance optimized for byte slice comparisons.
func NewBytesGolden(namespace string, opts ...Option[[]byte, []byte]) *Golden[[]byte, []byte] {
	return New(
		namespace,
		append(
			[]Option[[]byte, []byte]{
				WithSaver[[]byte, []byte](SaverFunc[[]byte](SaveBytes)),
				WithLoader[[]byte](LoaderFunc[[]byte](LoadBytes)),
			},
			opts...,
		)...,
	)
}

// NewStringGolden creates a Golden instance optimized for string comparisons.
func NewStringGolden(namespace string, opts ...Option[string, string]) *Golden[string, string] {
	return New(
		namespace,
		append(
			[]Option[string, string]{
				WithSaver[string, string](SaverFunc[string](SaveAsString[string])),
				WithLoader[string](LoaderFunc[string](LoadString)),
			},
			opts...,
		)...,
	)
}

// NewJSONGolden creates a Golden instance that uses JSON serialization for any type.
func NewJSONGolden[TActual, TExpected any](namespace string, opts ...Option[TActual, TExpected]) *Golden[TActual, TExpected] {
	return New(
		namespace,
		opts...,
	)
}

// WithGoldenFile configures a Golden instance to use a custom ReadWriter implementation.
func WithGoldenFile[TActual, TExpected any](goldenFile ReadWriter) Option[TActual, TExpected] {
	return func(g *Golden[TActual, TExpected]) {
		g.goldenFile = goldenFile
	}
}

// WithIsUpdater configures a Golden instance to use a custom update strategy.
func WithIsUpdater[TActual, TExpected any](updater IsUpdater) Option[TActual, TExpected] {
	return func(g *Golden[TActual, TExpected]) {
		g.updater = updater
	}
}

// WithSaver configures a Golden instance to use a custom saving strategy.
func WithSaver[TActual, TExpected any](saver Saver[TActual]) Option[TActual, TExpected] {
	return func(g *Golden[TActual, TExpected]) {
		g.saver = saver
	}
}

// WithLoader configures a Golden instance to use a custom loading strategy.
func WithLoader[TActual, TExpected any](loader Loader[TExpected]) Option[TActual, TExpected] {
	return func(g *Golden[TActual, TExpected]) {
		g.loader = loader
	}
}

// GetGoldenFile returns the underlying ReadWriter used by this Golden instance.
func (g *Golden[TActual, TExpected]) GetGoldenFile() ReadWriter {
	return g.goldenFile
}

// CleanupOption configures cleanup behavior.
type CleanupOption func(*CleanOptions)

// WithBackup enables backup creation during cleanup.
func WithBackup() CleanupOption {
	return func(opts *CleanOptions) {
		opts.Backup = true
	}
}

// WithDryRun enables dry-run mode (show what would be cleaned without deleting).
func WithDryRun() CleanupOption {
	return func(opts *CleanOptions) {
		opts.DryRun = true
	}
}

// WithPatterns sets file patterns to exclude from cleanup.
func WithPatterns(patterns ...string) CleanupOption {
	return func(opts *CleanOptions) {
		opts.Patterns = patterns
	}
}

// SetupCleanup configures cleanup for a Golden instance at test completion.
// This should be called at the test function level, not in subtests, to avoid cleanup timing issues.
//
// Cleanup behavior:
//   - If command line flags (-clean, -clean-dry-run) are set, those settings are used
//   - If custom options are provided, cleanup is enabled with those options
//   - If neither flags nor options are provided, no cleanup is performed
//
// Usage:
//
//	gilt.SetupCleanup(t, golden)                              // Only cleanup if -clean flags are used
//	gilt.SetupCleanup(t, golden, gilt.WithBackup())           // Force enable with backup
//	gilt.SetupCleanup(t, golden, gilt.WithDryRun())           // Force enable with dry-run
//	gilt.SetupCleanup(t, golden, gilt.WithPatterns("*.keep")) // Force enable with exclusions
func SetupCleanup[TActual, TExpected any](t *testing.T, g *Golden[TActual, TExpected], opts ...CleanupOption) {
	t.Helper()

	// Start with command line flags or defaults
	var options *CleanOptions
	if flagClean || flagCleanDryRun {
		options = &CleanOptions{
			DryRun: flagCleanDryRun,
			Backup: flagCleanBackup,
		}
	} else if len(opts) > 0 {
		// Only enable if custom options are provided
		options = &CleanOptions{}
	} else {
		// No cleanup if no flags and no custom options
		return
	}

	// Apply custom options
	for _, opt := range opts {
		opt(options)
	}

	// Register cleanup
	t.Cleanup(func() {
		if goldenFile, ok := g.goldenFile.(*GoldenFile); ok {
			CleanUnusedGoldenFiles(goldenFile, t, g.namespace, *options)
		}
	})
}
