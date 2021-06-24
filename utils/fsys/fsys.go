package fsys

import (
	"io/fs"
	"os"
)

type FS interface {
	GetAllFiles(dir string) ([]fs.FileInfo, error)
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
