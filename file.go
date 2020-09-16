package files

import (
	"os"
	"regexp"

	stacktrace "github.com/palantir/stacktrace"
)

// [TODO] add flock

// FileSize returns file size for the given path.
func FileSize(path string) (int64, error) {
	var err error
	fi, err := os.Stat(path)
	if err != nil {
		err = stacktrace.Propagate(err, "could not get file size at '%s'", path)
		return -1, err
	}
	if fi.IsDir() {
		err = stacktrace.NewError("'%s' is a directory", path)
		return -1, err
	}
	return fi.Size(), nil
}

// IsTemporaryFileName returns true if fn matches temporary file name pattern
func IsTemporaryFileName(fn string) bool {
	tmpFileNameRe := regexp.MustCompile(`\.tmp\.\d+$`)
	return tmpFileNameRe.MatchString(fn)
}
