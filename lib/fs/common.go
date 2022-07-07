package fs

import (
	"io"
	"os"
	"path/filepath"
)

const (
	WIN_OS_PROGDATA = "C:\\Temp\\postgres.backup"
)

func GetSize(path string) (int64, error) {
	if file, err := os.Open(path); err == nil {
		defer file.Close()
		info, err := file.Stat()
		if err != nil {
			return 0, err
		}
		if !info.IsDir() {
			return info.Size(), err
		}
	} else {
		return 0, err
	}
	return getDirectorySize(path)
}

func isDir(path string) (rs bool, err error) {
	if file, err := os.Open(path); err == nil {
		defer file.Close()
		info, err := file.Stat()
		if err != nil {
			return rs, err
		}
		return info.IsDir(), err
	} else {
		return rs, err
	}
}

func getDirectorySize(path string) (totalSize int64, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	return
}

func Remove(path string) error {
	return os.RemoveAll(path)
}

func TempDir() (string, error) {
	err := createIfNotExists(WIN_OS_PROGDATA)
	if err != nil {
		return "", err
	}
	return WIN_OS_PROGDATA, nil
}

func createIfNotExists(path string) error {
	return os.MkdirAll(path, 0777)
}

func CopyFile(src string, dst string) error {

	fdr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fdr.Close()

	fdw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fdw.Close()

	if _, err := io.Copy(fdw, fdr); err != nil {
		return err
	}

	return nil
}
