package earthfile2llb

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func hashDir(dirPath string) {
	filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("got %s\n", path)
		return nil
	})
}
