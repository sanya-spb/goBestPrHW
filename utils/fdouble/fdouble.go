package fdouble

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FDescr struct {
	path string
	hash string
	size uint64
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
func ScanDir(dir string, n *sync.WaitGroup, fileInfo chan<- FDescr) {
	defer n.Done()
	if cancelled() {
		return
	}
	for _, entry := range readDir(dir) {
		if entry.IsDir() {
			// fmt.Fprintf(os.Stdout, "DIR: %s\n", entry.Name())
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go ScanDir(subdir, n, fileInfo)
		} else {
			fileInfo <- FDescr{
				path: dir + string(filepath.Separator) + entry.Name(),
				hash: func(fPath string) string {
					content, err := os.Open(fPath)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%v\n", err)
					}
					hash := sha256.New()
					if _, err := io.Copy(hash, content); err != nil {
						fmt.Fprintf(os.Stderr, "%v\n", err)
					}
					return fmt.Sprintf("%x", hash.Sum(nil))
				}(dir + string(filepath.Separator) + entry.Name()),
				size: uint64(entry.Size()),
			}
		}
	}

}

var sema = make(chan struct{}, 20) // concurrency-limiting counting semaphore

// readDir returns the entries of directory dir.
func readDir(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: // acquire token
	case <-Done:
		return nil // cancelled
	}
	defer func() { <-sema }() // release token

	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		// Don't return: Readdir may return partial results.
	}
	return entries
}
