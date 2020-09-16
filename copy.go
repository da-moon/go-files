package files

import (
	"io/ioutil"
	"os"
	"path/filepath"

	primitives "github.com/da-moon/go-primitives"
	"github.com/palantir/stacktrace"
)

// CopyDir copies all directories, subdirectories and files recursively
func CopyDir(src, dest string) error {
	var err error
	src, err = filepath.Abs(src)
	if err != nil {
		err = stacktrace.Propagate(err, "could not copy src directory '%s' to destination '%s'", src, dest)
		return err
	}

	dest, err = filepath.Abs(dest)
	if err != nil {
		err = stacktrace.Propagate(err, "could not copy src directory '%s' to destination '%s'", src, dest)
		return err
	}

	files, err := ReadDirFiles(src, "")
	if err != nil {
		err = stacktrace.Propagate(err, "could not copy src directory '%s' to destination '%s'", src, dest)
		return err
	}

	for _, f := range files {
		dp := primitives.PathJoin(dest, f)
		sp := primitives.PathJoin(src, f)
		err = os.MkdirAll(filepath.Dir(dp), 0777)
		if err != nil {
			err = stacktrace.Propagate(err, "could not copy src directory '%s' to destination '%s'", src, dest)
			return err
		}
		err := CopyFile(sp, dp)
		if err != nil {
			err = stacktrace.Propagate(err, "could not copy src directory '%s' to destination '%s'", src, dest)
			return err
		}
	}
	return nil
}

// CopyFile copies a file from src to destination
func CopyFile(src, dest string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		err = stacktrace.Propagate(err, "could not copy src file '%s' to destination '%s'", src, dest)
		return err
	}

	err = ioutil.WriteFile(dest, data, 0666)
	if err != nil {
		err = stacktrace.Propagate(err, "could not copy src file '%s' to destination '%s'", src, dest)
		return err
	}
	return nil
}
