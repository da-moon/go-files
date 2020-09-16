package files

import (
	"os"
	"path/filepath"

	primitives "github.com/da-moon/go-primitives"
	"github.com/palantir/stacktrace"
)

// HardLinkFiles makes hard links for all the files from src in dest.
func HardLinkFiles(src, dest string) error {
	var err error
	err = mkdirSync(dest)
	if err != nil {
		err = stacktrace.Propagate(err, "cannot create dest '%s'", dest)
		return err
	}

	d, _, err := OpenPath(src)
	if err != nil {
		err = stacktrace.Propagate(err, "cannot open src '%s'", src)
		return err
	}
	// [TODO] maybe, its a good idea to check for errors and log closing of src
	defer d.Close()
	fis, err := d.Readdir(-1)
	if err != nil {
		err = stacktrace.Propagate(err, "cannot read files in scrDir=%q: %w", src, err)
		return err
	}
	for _, fi := range fis {
		if IsDirOrSymlink(fi) {
			continue
		}
		fn := fi.Name()
		srcPath := primitives.PathJoin(src, fn)
		dstPath := primitives.PathJoin(dest, fn)
		err := os.Link(srcPath, dstPath)
		if err != nil {
			return err
		}
	}
	err = SyncPath(dest)
	if err != nil {
		err = stacktrace.Propagate(err, "could not sync '%s'", dest)
		return err
	}
	return nil
}

// IsDirOrSymlink returns true if fi is directory or symlink.
func IsDirOrSymlink(fi os.FileInfo) bool {
	return fi.IsDir() || (fi.Mode()&os.ModeSymlink == os.ModeSymlink)
}

// SymlinkRelative creates relative symlink for srcPath in dstPath.
func SymlinkRelative(srcPath, dstPath string) error {
	baseDir := filepath.Dir(dstPath)
	srcPathRel, err := filepath.Rel(baseDir, srcPath)
	if err != nil {
		return stacktrace.Propagate(err, "cannot make relative path for srcPath=%q: %w", srcPath, err)
	}
	err = os.Symlink(srcPathRel, dstPath)
	if err != nil {
		return err
	}
	return nil
}
