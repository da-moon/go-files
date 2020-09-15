package files

import (
	"os"
	"path/filepath"

	stacktrace "github.com/palantir/stacktrace"
)

// FindFile searches a root directory recursively for files with a pattern
func FindFile(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			err = stacktrace.Propagate(err, "filepath did not match pattern")
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		err = stacktrace.Propagate(err, "could not find file with root path '%s' and pattern '%s'", root, pattern)
		return nil, err
	}
	return matches, nil
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
	if fi.Size() == 0 {
		os.Remove(path)
		err = stacktrace.NewError("file at '%v' was empty", path)
		return nil, nil, err
	}

	return f, fi, nil
}
