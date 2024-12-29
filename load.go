package gilt

import (
	"bufio"
	"encoding/json"
	"iter"
	"testing"
)

var _ Loader[any] = (*LoaderFunc[any])(nil)

type LoaderFunc[TExpected any] func(t *testing.T, name string, reader Reader) TExpected

func (f LoaderFunc[TExpected]) Load(t *testing.T, name string, reader Reader) TExpected {
	t.Helper()
	return f(t, name, reader)
}

func LoadBytes(t *testing.T, name string, reader Reader) []byte {
	t.Helper()
	return reader.Read(t, name)
}

func LoadAndUnmarshalJSON[TExpected any](t *testing.T, name string, reader Reader) TExpected {
	t.Helper()
	b := reader.Read(t, name)
	var expected TExpected
	if err := json.Unmarshal(b, &expected); err != nil {
		t.Fatalf("failed to unmarshal golden file %s: %v", name, err)
	}
	return expected
}

func LoadString(t *testing.T, name string, reader Reader) string {
	t.Helper()
	b := reader.Read(t, name)
	return string(b)
}

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
