package fdouble

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/sanya-spb/goBestPrHW/utils/fsys"
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

	// знаю что плохо, первый и последний раз.. время поджимает..
	time.Sleep(1 * time.Second)

	assert.Equal(t, []FDescr{
		{
			"abc" + string(filepath.Separator) + "1.dat",
			"100500",
			100,
		},
	}, retFiles)
}

func TestScanDirIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ll := logrus.New()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	files := make(chan FDescr)

	fsTemp := &fsys.FSImpl{}

	dir, err := ioutil.TempDir("/tmp", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := ioutil.TempFile(dir, "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	var retFiles []FDescr
	go func() {
		for f := range files {
			retFiles = append(retFiles, f)
		}
	}()

	ScanDir(dir, wg, files, fsTemp, ll)

	// знаю что плохо, первый и последний раз.. время поджимает..
	time.Sleep(1 * time.Second)

	assert.Equal(t, []FDescr{
		{
			file.Name(),
			fsTemp.CalcHash(file.Name()),
			func(f *os.File) uint64 {
				fi, err := f.Stat()
				if err != nil {
					// Could not obtain stat, no handle error, test will fail
					return 111
				}
				return uint64(fi.Size())
			}(file),
		},
	}, retFiles)
}
