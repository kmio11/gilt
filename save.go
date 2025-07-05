package gilt

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

var _ Saver[any] = (*SaverFunc[any])(nil)

// SaverFunc is a function type that implements the Saver interface.
type SaverFunc[TActual any] func(t *testing.T, actual TActual, name string, writer Writer)

// Save implements the Saver interface by calling the underlying function.
func (f SaverFunc[TActual]) Save(t *testing.T, actual TActual, name string, writer Writer) {
	t.Helper()
	f(t, actual, name, writer)
}

// SaveBytes saves byte data directly to a golden file.
func SaveBytes(t *testing.T, actual []byte, name string, writer Writer) {
	t.Helper()
	writer.Write(t, name, actual)
}

// SaveAsJSON saves data as JSON with indentation to a golden file.
func SaveAsJSON[TActual any](t *testing.T, actual TActual, name string, writer Writer) {
	t.Helper()
	b, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writer.Write(t, name, b)
}

// SaveAsString saves data as a string representation to a golden file.
func SaveAsString[TActual any](t *testing.T, actual TActual, name string, writer Writer) {
	t.Helper()
	b := []byte(fmt.Sprintf("%v", actual))
	writer.Write(t, name, b)
}

// SaveLines saves a slice of string-like elements as line-separated content to a golden file.
func SaveLines[TActualElm ~string](t *testing.T, actual []TActualElm, name string, writer Writer) {
	t.Helper()
	strArr := make([]string, len(actual))
	for i, elm := range actual {
		strArr[i] = fmt.Sprintf("%v", elm)
	}
	writer.Write(t, name, []byte(strings.Join(strArr, "\n")))
}
