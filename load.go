package gilt

import (
	"bufio"
	"encoding/json"
	"iter"
	"testing"
)

var _ Loader[any] = (*LoaderFunc[any])(nil)

// LoaderFunc is a function type that implements the Loader interface.
type LoaderFunc[TExpected any] func(t *testing.T, name string, reader Reader) TExpected

// Load implements the Loader interface by calling the underlying function.
func (f LoaderFunc[TExpected]) Load(t *testing.T, name string, reader Reader) TExpected {
	t.Helper()
	return f(t, name, reader)
}

// LoadBytes loads raw byte data from a golden file.
func LoadBytes(t *testing.T, name string, reader Reader) []byte {
	t.Helper()
	return reader.Read(t, name)
}

// LoadAndUnmarshalJSON loads and unmarshals JSON data from a golden file into the specified type.
func LoadAndUnmarshalJSON[TExpected any](t *testing.T, name string, reader Reader) TExpected {
	t.Helper()
	b := reader.Read(t, name)
	var expected TExpected
	if err := json.Unmarshal(b, &expected); err != nil {
		t.Fatalf("failed to unmarshal golden file %s: %v", name, err)
	}
	return expected
}

// LoadString loads the contents of a golden file as a string.
func LoadString(t *testing.T, name string, reader Reader) string {
	t.Helper()
	b := reader.Read(t, name)
	return string(b)
}

// LoadLines loads the contents of a golden file as an iterator over lines.
func LoadLines(t *testing.T, name string, reader Reader) iter.Seq[string] {
	t.Helper()

	o, ok := reader.(Opener)
	if !ok {
		t.Fatalf("reader does not implement Opener")
	}

	return func(yield func(string) bool) {
		f := o.Open(t, name)
		t.Cleanup(func() { f.Close() })

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}
}
