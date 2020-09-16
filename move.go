package files

import (
	"os"
	"path/filepath"

	"github.com/palantir/stacktrace"
)

// [TODO] add move function

// Rename safely renames a file.
func Rename(from, to string) error {
	var err error
	err = os.Rename(from, to)
	if err != nil {
		err = stacktrace.Propagate(err, "could not rename '%s' to '%s'", from, to)
		return err
	}
	err = SyncPath(filepath.Dir(to))
	if err != nil {
		err = stacktrace.Propagate(err, "could not rename '%s' to '%s'", from, to)
		return err
	}
	return nil
}
