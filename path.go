package files

import (
	"os"

	"github.com/palantir/stacktrace"
)

// [TODO] add direct io opener

// PathExist returns whether the given path exists.
func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// OpenPath ...
func OpenPath(path string) (*os.File, os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		err = stacktrace.Propagate(err, "Error reading '%s'", path)
		return nil, nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		err = stacktrace.Propagate(err, "Error reading '%s'", path)
		return nil, nil, err
	}

	return f, fi, nil
}

// SafeOpenPath removes empty files after opening
// empty files are most often result of a failed io
func SafeOpenPath(path string) (*os.File, os.FileInfo, error) {
	f, fi, err := OpenPath(path)
	if err != nil {
		err = stacktrace.Propagate(err, "could not open '%s'", path)
		return nil, nil, err
	}
	if !fi.IsDir() && fi.Size() == 0 {
		err = os.Remove(path)
		if err != nil {
			err = stacktrace.Propagate(err, "could not open '%s'", path)
			return nil, nil, err
		}
	}
	return f, fi, err
}

// SyncPath makes sure file at a certain path is synced with physical disk
func SyncPath(path string) error {
	var err error
	f, _, err := OpenPath(path)
	if err != nil {
		err = stacktrace.Propagate(err, "could not sync path '%s'", path)
		return err
	}

	err = f.Sync()
	if err != nil {
		_ = f.Close()
		err = stacktrace.Propagate(err, "could not flush path '%s' to disk", path)
		return err
	}
	err = f.Close()
	if err != nil {
		err = stacktrace.Propagate(err, "could not close file descriptor at '%s' to disk", path)
		return err
	}
	return nil
}
