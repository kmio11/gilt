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

// SetIsUpdateFunc sets a custom function to determine when golden files should be updated.
func SetIsUpdateFunc(f func(t *testing.T, name string) bool) {
	isUpdateFunc = f
}

var _ IsUpdater = (*IsUpdaterFunc)(nil)

// IsUpdaterFunc is a function type that implements the IsUpdater interface.
type IsUpdaterFunc func(t *testing.T, name string) bool

// IsUpdate implements the IsUpdater interface by calling the underlying function.
func (f IsUpdaterFunc) IsUpdate(t *testing.T, name string) bool {
	return f(t, name)
}
