package gilt

import (
	"os"
	"testing"
)

type (
	ReadWriter interface {
		Reader
		Writer
	}
	Reader interface {
		Read(t *testing.T, name string) []byte
	}
	Writer interface {
		Write(t *testing.T, name string, data []byte)
	}
	Opener interface {
		Open(t *testing.T, name string) *os.File
	}
)

type IsUpdater interface {
	IsUpdate(t *testing.T, name string) bool
}

type Saver[TActual any] interface {
	Save(t *testing.T, actual TActual, name string, writer Writer)
}

type Loader[TExpected any] interface {
	Load(t *testing.T, name string, reader Reader) TExpected
}

type Golden[TActual, TExpected any] struct {
	updater    IsUpdater
	saver      Saver[TActual]
	loader     Loader[TExpected]
	goldenFile ReadWriter
}

type Option[TActual, TExpected any] func(*Golden[TActual, TExpected])

func New[TActual any, TExpected any](
	namespace string,
	opts ...Option[TActual, TExpected],
) *Golden[TActual, TExpected] {
	g := &Golden[TActual, TExpected]{
		goldenFile: NewGoldenFile(namespace),
		updater:    IsUpdaterFunc(isUpdateFunc),
		saver:      SaverFunc[TActual](SaveAsJSON[TActual]),
		loader:     LoaderFunc[TExpected](LoadAndUnmarshalJSON[TExpected]),
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (g *Golden[TActual, TExpected]) Assert(
	t *testing.T, actual TActual, name string,
	assertion func(t *testing.T, actual TActual, expected TExpected),
) {
	t.Helper()

	if g.updater.IsUpdate(t, name) {
		g.saver.Save(t, actual, name, g.goldenFile)
		return
	}

	expectedData := g.loader.Load(t, name, g.goldenFile)
	assertion(t, actual, expectedData)
}

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

func NewJSONGolden[TActual, TExpected any](namespace string, opts ...Option[TActual, TExpected]) *Golden[TActual, TExpected] {
	return New(
		namespace,
		opts...,
	)
}

func WithGoldenFile[TActual, TExpected any](goldenFile ReadWriter) Option[TActual, TExpected] {
	return func(g *Golden[TActual, TExpected]) {
		g.goldenFile = goldenFile
	}
}

func WithIsUpdater[TActual, TExpected any](updater IsUpdater) Option[TActual, TExpected] {
	return func(g *Golden[TActual, TExpected]) {
		g.updater = updater
	}
}

func WithSaver[TActual, TExpected any](saver Saver[TActual]) Option[TActual, TExpected] {
	return func(g *Golden[TActual, TExpected]) {
		g.saver = saver
	}
}

func WithLoader[TActual, TExpected any](loader Loader[TExpected]) Option[TActual, TExpected] {
	return func(g *Golden[TActual, TExpected]) {
		g.loader = loader
	}
}
