package gilt_test

import (
	"fmt"
	"iter"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kmio11/gilt"
)

func hello(to string) string {
	return "Hello, " + to + "!"
}

func TestNewStringGolden(t *testing.T) {
	golden := gilt.NewStringGolden(t.Name())

	tests := []struct {
		name string
	}{
		{"world"},
		{"gopher"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := hello(tt.name)
			golden.Assert(t, actual, tt.name, func(t *testing.T, actual string, expected string) {
				if actual != expected {
					t.Errorf("expected: %s\n, but was: %s", expected, actual)
				}
			})
		})
	}
}

func TestNewBytesGolden(t *testing.T) {
	golden := gilt.NewBytesGolden(t.Name())

	tests := []struct {
		name string
	}{
		{"world"},
		{"gopher"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := []byte(hello(tt.name))
			golden.Assert(t, actual, tt.name, func(t *testing.T, actual []byte, expected []byte) {
				if string(actual) != string(expected) {
					t.Errorf("expected: %s\n, but was: %s", expected, actual)
				}
			})
		})
	}
}

type HelloMessage struct {
	Message string `json:"message"`
}

func helloMessage(to string) HelloMessage {
	return HelloMessage{"Hello, " + to + "!"}
}

func TestNewJSONGolden(t *testing.T) {
	golden := gilt.NewJSONGolden[HelloMessage, HelloMessage](t.Name())

	tests := []struct {
		name string
	}{
		{"world"},
		{"gopher"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := helloMessage(tt.name)
			golden.Assert(t, actual, tt.name, func(t *testing.T, actual HelloMessage, expected HelloMessage) {
				if actual.Message != expected.Message {
					t.Errorf("expected: %v\n, but was: %v", expected, actual)
				}
			})
		})
	}
}

func TestNewJSONGolden_with_custom_filepath(t *testing.T) {
	golden := gilt.NewJSONGolden(
		t.Name(),
		gilt.WithGoldenFile[HelloMessage, HelloMessage](
			gilt.NewGoldenFile(
				t.Name(),
				gilt.GoldenFileWithPathHandler(func(t *testing.T, namespace string, name string) string {
					// save the golden files with .json extension
					return filepath.Join("testdata", namespace, "golden", fmt.Sprintf("%s.json", name))
				}),
			),
		),
	)

	tests := []struct {
		name string
	}{
		{"world"},
		{"gopher"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := helloMessage(tt.name)
			golden.Assert(t, actual, tt.name, func(t *testing.T, actual HelloMessage, expected HelloMessage) {
				if actual.Message != expected.Message {
					t.Errorf("expected: %v\n, but was: %v", expected, actual)
				}
			})
		})
	}
}

func hellos(to []string) []string {
	var messages []string
	for _, name := range to {
		messages = append(messages, hello(name))
	}
	return messages

}

func TestNew(t *testing.T) {
	golden := gilt.New(
		t.Name(),
		gilt.WithSaver[[]string, iter.Seq[string]](
			gilt.SaverFunc[[]string](gilt.SaveLines[string]),
		),
		gilt.WithLoader[[]string](
			gilt.LoaderFunc[iter.Seq[string]](gilt.LoadLines),
		),
	)

	tests := []struct {
		names []string
	}{
		{[]string{"world", "gopher"}},
	}

	for _, tt := range tests {
		testName := strings.Join(tt.names, "_")
		t.Run(testName, func(t *testing.T) {
			actual := hellos(tt.names)
			golden.Assert(t, actual, testName, func(t *testing.T, actual []string, expected iter.Seq[string]) {
				line := 0
				for elm := range expected {
					if actual[line] != elm {
						t.Errorf("(line %d) expected: %s\n, but was: %s", line, elm, actual[line])
					}
					line++
				}
			})
		})
	}
}
