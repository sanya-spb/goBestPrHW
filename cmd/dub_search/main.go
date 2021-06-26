package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sanya-spb/goBestPrHW/internal/config"
	"github.com/sanya-spb/goBestPrHW/internal/fdouble"
	"github.com/sanya-spb/goBestPrHW/pkg/version"

	"github.com/sirupsen/logrus"
)

type APP struct {
	Version version.AppVersion
	Config  config.Config
	logrus  *logrus.Logger
}

var MyApp *APP = new(APP)

func main() {
	MyApp.logrus = logrus.New()
	MyApp.logrus.Formatter = new(logrus.TextFormatter)
	MyApp.logrus.Formatter.(*logrus.TextFormatter).DisableTimestamp = true
	MyApp.logrus.Out = os.Stdout

	MyApp.Version = *version.Version
	MyApp.Config = *config.NewConfig()

	if MyApp.Config.Debug {
		MyApp.logrus.Level = logrus.TraceLevel
	} else {
		MyApp.logrus.Level = logrus.InfoLevel
	}

	// fmt.Printf("version: %+v\n", MyApp.Version)
	MyApp.logrus.WithFields(logrus.Fields{
		"version":    MyApp.Version.Version,
		"commit":     MyApp.Version.Commit,
		"build time": MyApp.Version.BuildTime,
		"copyright":  MyApp.Version.Copyright,
	}).Debug("version")

	// fmt.Printf("config: %+v\n", MyApp.Config)
	MyApp.logrus.WithFields(logrus.Fields{
		"dirs":     MyApp.Config.Dirs,
		"dfactor>": MyApp.Config.DFactor,
	}).Debug("config")

	// Cancel traversal when input is detected.
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(fdouble.Done)
	}()

	// Traverse each root of the file tree in parallel.
	fileHashesCh := make(chan fdouble.FDescr)
	// subDirCh := make(chan fdouble.SubDir)
	var n sync.WaitGroup
	for _, root := range MyApp.Config.Dirs {
		n.Add(1)
		go fdouble.ScanDir(strings.TrimRight(root, string(filepath.Separator)), &n, fileHashesCh, MyApp.logrus)
	}
	go func() {
		n.Wait()
		close(fileHashesCh)
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
			for range fileHashesCh {
				// Do nothing.
			}
			fmt.Printf("Cancelled by user\n")
			os.Exit(1)
		case fHash, ok := <-fileHashesCh:
			if !ok {
				break loop // fileHashes was closed
			}
			// using map features for solve this task:
			if fHash.Hash() != "" {
				key := keyDBFileHash{
					hash: fHash.Hash(),
					size: fHash.Size(),
				}
				DBfileHash[key] = append(DBfileHash[key], fHash.Path())
			}
		}

	}

	for k, vPath := range DBfileHash {
		if len(vPath) > int(MyApp.Config.DFactor) {
			// fmt.Printf("hash: %s, size: %d\n", k.hash, k.size)
			// fmt.Printf("found %d doubles:\n", len(vPath))
			for ii, v := range vPath {
				// fmt.Printf("  %s\n", v)
				MyApp.logrus.WithFields(logrus.Fields{
					"hash":    k.hash,
					"size":    k.size,
					"doubles": len(vPath),
					"id":      ii,
				}).Info(v)
			}
		}
	}
}
