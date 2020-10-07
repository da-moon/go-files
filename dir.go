package files

import (
	"os"
	"path/filepath"
	"strings"

	stacktrace "github.com/palantir/stacktrace"
)

// DirSize returns size of a target directory
func DirSize(dir string) (int64, error) {
	var size int64
	err := filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		size += info.Size()
		return nil
	})
	return size, err
}

// ReadDirFiles searches a root directory recursively for files with a pattern
func ReadDirFiles(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if pattern != "" {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				err = stacktrace.Propagate(err, "filepath did not match pattern")
				return err
			} else if matched {
				relativePath := strings.TrimPrefix(path, root)
				if len(relativePath) > 0 {
					matches = append(matches, relativePath)
				}
			}
			return nil
		}
		relativePath := strings.TrimPrefix(path, root)
		if len(relativePath) > 0 {
			matches = append(matches, relativePath)
		}
		return nil
	})
	if err != nil {
		err = stacktrace.Propagate(err, "could not find file with root path '%s' and pattern '%s'", root, pattern)
		return nil, err
	}
	return matches, nil
}

// SafeMkdirAll creates directory tree in case it doesn't exist
// if it exists, it would fail
func SafeMkdirAll(path string) error {
	var err error
	if PathExist(path) {
		err = stacktrace.NewError("'%s' already exists", path)
		return err
	}
	err = mkdirSync(path)
	if err != nil {
		err = stacktrace.Propagate(err, "could not safely create directory tree at '%s' and sync it with disk", path)
		return err
	}
	return nil
}

// MkdirAll creates directory tree and won't return error if it exists
func MkdirAll(path string) error {
	var err error
	if PathExist(path) {
		return nil
	}
	err = mkdirSync(path)
	if err != nil {
		err = stacktrace.Propagate(err, "could not safely create directory tree at '%s' and sync it with disk", path)
		return err
	}
	return nil
}
func mkdirSync(path string) error {
	var err error
	err = os.MkdirAll(path, 0755)
	if err != nil {
		err = stacktrace.Propagate(err, "could not create directory at '%s'", path)
		return err
	}
	err = SyncPath(filepath.Dir(path))
	if err != nil {
		err = stacktrace.Propagate(err, "could not create directory at '%s'", path)
		return err
	}
	return nil
}
