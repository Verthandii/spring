package utils

import (
	"bytes"
	"encoding/csv"
)

// NewCSV 数组转 CSV 文件流
func NewCSV(records [][]string) ([]byte, error) {
	buf := bytes.NewBuffer([]byte("\xEF\xBB\xBF"))
	w := csv.NewWriter(buf)
	if err := w.WriteAll(records); err != nil {
		return nil, err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
