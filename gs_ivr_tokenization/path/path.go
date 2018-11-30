package path

import (
	"os"
	"path/filepath"
)

// RelativePath returns the absolte path from relative path
func RelativePath(path string) string {
	if len(path) > 2 && path[:2] == "./" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		path = dir + path[1:]
	}
	return path
}
