package gilt

import (
	"flag"
	"testing"
)

var flagUpdate bool

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

// isUpdateFunc is a function variable that determines whether an update should be performed.
var isUpdateFunc = func(t *testing.T, name string) bool {
	return flagUpdate
}

func SetIsUpdateFunc(f func(t *testing.T, name string) bool) {
	isUpdateFunc = f
}

var _ IsUpdater = (*IsUpdaterFunc)(nil)

type IsUpdaterFunc func(t *testing.T, name string) bool

func (f IsUpdaterFunc) IsUpdate(t *testing.T, name string) bool {
	return f(t, name)
}
