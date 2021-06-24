package fdouble

import (
	"io/fs"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/sanya-spb/goBestPrHW/utils/fsys/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type FileInfoStub struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (f *FileInfoStub) Name() string {
	return f.name
}
func (f *FileInfoStub) Size() int64 {
	return f.size
}
func (f *FileInfoStub) Mode() fs.FileMode {
	return f.mode
}
func (f *FileInfoStub) ModTime() time.Time {
	return f.modTime
}
func (f *FileInfoStub) IsDir() bool {
	return f.isDir
}
func (f *FileInfoStub) Sys() interface{} {
	return f.sys
}

func TestScanDir(t *testing.T) {
	ll := logrus.New()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	files := make(chan FDescr)

	// fs := &FSStub{}
	fsMock := &mocks.FS{}

	var fsTest []fs.FileInfo
	fsTest = append(fsTest, &FileInfoStub{
		name:  "1.dat",
		size:  100,
		mode:  0,
		isDir: false,
	})

	// fsMock.On("GetAllFiles", mock.Anything).Return(fsTest, nil)
	// fsMock.On("CalcHash", mock.Anything).Return("100500")
	fsMock.On("GetAllFiles", string("abc")).Return(fsTest, nil)
	fsMock.On("CalcHash", string("abc/1.dat")).Return("100500")

	var retFiles []FDescr
	go func() {
		for f := range files {
			retFiles = append(retFiles, f)
		}
	}()

	ScanDir("abc", wg, files, fsMock, ll)

	time.Sleep(1 * time.Second)

	assert.Equal(t, []FDescr{
		{
			"abc" + string(filepath.Separator) + "1.dat",
			"100500",
			100,
		},
	}, retFiles)
}
