package utils

import "os"

// PathExists 文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Remove 删除文件,为了defer用
func Remove(file string) {
	_ = os.Remove(file)
}
