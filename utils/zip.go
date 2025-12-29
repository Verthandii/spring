package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateRawZipFile(zipFileName string, files []string) (string, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	zipDir, err := os.MkdirTemp("", "zip")
	if err != nil {
		return "", err
	}

	zipFileName = filepath.Join(zipDir, zipFileName)
	archive, err := os.OpenFile(zipFileName, flags, 0644)
	if err != nil {
		return "", err
	}
	defer archive.Close()

	zipw := zip.NewWriter(archive)
	defer zipw.Close()

	for _, f := range files {
		if err = appendFiles(".", f, zipw); err != nil {
			return "", err
		}
	}
	return zipFileName, nil
}

func CreateZipFile(zipFileName string, archives map[string][]string) (string, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	zipDir, err := os.MkdirTemp("", "zip")
	if err != nil {
		return "", err
	}

	zipFileName = filepath.Join(zipDir, zipFileName)
	archive, err := os.OpenFile(zipFileName, flags, 0644)
	if err != nil {
		return "", err
	}
	defer archive.Close()

	zipw := zip.NewWriter(archive)
	defer zipw.Close()

	for dir, fs := range archives {
		archiveDir := filepath.Join(zipDir, dir)

		if err = os.MkdirAll(archiveDir, os.ModePerm); err != nil {
			return "", err
		}
		for _, f := range fs {
			if err = appendFiles(dir, f, zipw); err != nil {
				return "", err
			}
		}
	}
	return zipFileName, nil
}

func appendFiles(dir, filename string, zipw *zip.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open %s: %s", filename, err)
	}
	defer file.Close()

	base := filepath.Base(filename)
	filename = filepath.Join(dir, base)
	wr, err := zipw.Create(filename)
	if err != nil {
		msg := "failed to create entry for %s in zip file: %s"
		return fmt.Errorf(msg, filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("failed to write %s to zip: %s", filename, err)
	}
	return nil
}
