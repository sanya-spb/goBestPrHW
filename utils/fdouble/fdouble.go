package fdouble

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

type FDescr struct {
	path string
	hash string
	size uint64
}

type SubDir struct {
	pathFrom string
	pathTo   string
}

func (f *FDescr) Path() string {
	return f.path
}

func (f *FDescr) Hash() string {
	return f.hash
}

func (f *FDescr) Size() uint64 {
	return f.size
}

func (f *SubDir) PathFrom() string {
	return f.pathFrom
}

func (f *SubDir) PathTo() string {
	return f.pathTo
}

var Done = make(chan struct{})

func cancelled() bool {
	select {
	case <-Done:
		return true
	default:
		return false
	}
}

// ScanDir recursively walks the file tree and sends the fDescr of each found file.
func ScanDir(dir string, n *sync.WaitGroup, fileInfo chan<- FDescr, ll *logrus.Logger) {
	defer n.Done()

	// hook for panic in task 4
	defer func() {
		if err := recover(); err != nil {
			switch x := err.(type) {
			case error:
				ll.WithFields(logrus.Fields{
					"Unknown error": x.Error(),
				}).Error("sub Dir")
			case SubDir:
				ll.WithFields(logrus.Fields{
					"DirFrom": x.PathFrom(),
					"DirTo":   x.PathTo(),
				}).Error("sub Dir")
			}
		}
	}()

	if cancelled() {
		return
	}
	for _, entry := range readDir(dir, ll) {
		if entry.IsDir() {
			// fmt.Fprintf(os.Stdout, "DIR: %s\n", entry.Name())
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			ll.WithFields(logrus.Fields{
				"DirFrom": dir,
				"DirTo":   subdir,
			}).Debug("sub Dir")
			go ScanDir(subdir, n, fileInfo, ll)

			// Panic for task 4
			// panic(SubDir{pathFrom: dir, pathTo: subdir})
		} else if entry.Mode().IsRegular() {
			fileInfo <- FDescr{
				path: dir + string(filepath.Separator) + entry.Name(),
				hash: func(fPath string) string {
					content, err := os.Open(fPath)
					if err != nil {
						// fmt.Fprintf(os.Stderr, "%v\n", err)
						ll.WithFields(logrus.Fields{
							"fPath": fPath,
						}).Error(err.Error())
					} else {
						hash := sha256.New()
						if _, err := io.Copy(hash, content); err != nil {
							fmt.Fprintf(os.Stderr, "%v\n", err)
							ll.WithFields(logrus.Fields{
								"fPath": fPath,
							}).Error(err.Error())
						}
						return fmt.Sprintf("%x", hash.Sum(nil))
					}
					return ""
				}(dir + string(filepath.Separator) + entry.Name()),
				size: uint64(entry.Size()),
			}
		}
	}

}

var sema = make(chan struct{}, 20) // concurrency-limiting counting semaphore

// readDir returns the entries of directory dir.
func readDir(dir string, ll *logrus.Logger) []os.FileInfo {
	select {
	case sema <- struct{}{}: // acquire token
	case <-Done:
		return nil // cancelled
	}
	defer func() { <-sema }() // release token

	f, err := os.Open(dir)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "%v\n", err)
		ll.WithFields(logrus.Fields{
			"fPath": dir,
		}).Error(err.Error())
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {
		// fmt.Fprintf(os.Stderr, "%v\n", err)
		ll.WithFields(logrus.Fields{
			"fPath": dir,
		}).Error(err.Error())
		// Don't return: Readdir may return partial results.
	}
	return entries
}
