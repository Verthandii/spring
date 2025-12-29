package source

import (
	"io"
	"os"
)

type file struct {
	path string
}

// NewFile 创建一个文件的数据源
func NewFile(path string) Source {
	return &file{path: path}
}

// Load 加载文件数据
func (f *file) Load() ([]byte, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
