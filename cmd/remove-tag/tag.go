package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func generateFile(root string) error {
	return filepath.Walk(root, func(path string, f fs.FileInfo, err error) error {
		if !strings.HasSuffix(f.Name(), ".pb.go") {
			return nil
		}
		fname := path
		fbytes, err := os.ReadFile(fname)
		if err != nil {
			return err
		}
		fdata := strings.ReplaceAll(string(fbytes), ",omitempty", "")
		return os.WriteFile(fname, []byte(fdata), 0644)
	})
}
