package gilt

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

var _ Saver[any] = (*SaverFunc[any])(nil)

type SaverFunc[TActual any] func(t *testing.T, actual TActual, name string, writer Writer)

func (f SaverFunc[TActual]) Save(t *testing.T, actual TActual, name string, writer Writer) {
	t.Helper()
	f(t, actual, name, writer)
}

func SaveBytes(t *testing.T, actual []byte, name string, writer Writer) {
	t.Helper()
	writer.Write(t, name, actual)
}

func SaveAsJSON[TActual any](t *testing.T, actual TActual, name string, writer Writer) {
	t.Helper()
	b, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writer.Write(t, name, b)
}

func SaveAsString[TActual any](t *testing.T, actual TActual, name string, writer Writer) {
	t.Helper()
	b := []byte(fmt.Sprintf("%v", actual))
	writer.Write(t, name, b)
}

func SaveLines[TActualElm ~string](t *testing.T, actual []TActualElm, name string, writer Writer) {
	t.Helper()
	strArr := make([]string, len(actual))
	for i, elm := range actual {
		strArr[i] = fmt.Sprintf("%v", elm)
	}
	writer.Write(t, name, []byte(strings.Join(strArr, "\n")))
}
