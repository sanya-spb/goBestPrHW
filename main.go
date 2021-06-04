package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sanya-spb/goLev2HW/utils/config"
	"github.com/sanya-spb/goLev2HW/utils/fdouble"
	"github.com/sanya-spb/goLev2HW/utils/version"
)

type APP struct {
	Version version.AppVersion
	Config  config.Config
}

var MyApp *APP = new(APP)

func main() {
	MyApp.Version = *version.Version
	MyApp.Config = *config.NewConfig()

	fmt.Printf("version: %+v\n", MyApp.Version)
	fmt.Printf("config: %+v\n", MyApp.Config)

	// Cancel traversal when input is detected.
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(fdouble.Done)
	}()

	// Traverse each root of the file tree in parallel.
	fileHashes := make(chan fdouble.FDescr)
	var n sync.WaitGroup
	for _, root := range MyApp.Config.Dirs {
		n.Add(1)
		go fdouble.ScanDir(strings.TrimRight(root, string(filepath.Separator)), &n, fileHashes)
	}
	go func() {
		n.Wait()
		close(fileHashes)
	}()

	type keyDBFileHash struct {
		hash string
		size uint64
	}

	DBfileHash := map[keyDBFileHash][]string{}

	// var cancelled = false
loop:
	for {
		select {
		case <-fdouble.Done:
			// Drain fileHashes to allow existing goroutines to finish.
			for range fileHashes {
				// Do nothing.
			}
			fmt.Printf("Cancelled by user\n")
			os.Exit(1)
		case fHash, ok := <-fileHashes:
			if !ok {
				break loop // fileHashes was closed
			}
			// using map features for solve this task:
			key := keyDBFileHash{
				hash: fHash.Hash(),
				size: fHash.Size(),
			}
			DBfileHash[key] = append(DBfileHash[key], fHash.Path())
		}

	}

	for k, vPath := range DBfileHash {
		if len(vPath) > int(MyApp.Config.DFactor) {
			fmt.Printf("hash: %s, size: %d\n", k.hash, k.size)
			fmt.Printf("found %d doubles:\n", len(vPath))
			for _, v := range vPath {
				fmt.Printf("  %s\n", v)
			}
		}
	}
}
