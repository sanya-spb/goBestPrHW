package fsys

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
)

type FS interface {
	GetAllFiles(dir string) ([]fs.FileInfo, error)
	CalcHash(path string) string
}

type FSImpl struct {
}

func (fs FSImpl) GetAllFiles(dir string) ([]fs.FileInfo, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {

	}
	return entries, nil
}

func (fs FSImpl) CalcHash(path string) string {
	content, err := os.Open(path)
	if err != nil {

	} else {
		hash := sha256.New()
		if _, err := io.Copy(hash, content); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		return fmt.Sprintf("%x", hash.Sum(nil))
	}
	return ""
}
