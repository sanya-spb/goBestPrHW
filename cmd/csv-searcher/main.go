package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	lErr *log.Logger
	lOut *log.Logger
)

func main() {
	// Init our app
	app, err := newApp()
	if err != nil {
		// fmt.Fprintln(os.Stderr, err.Error())
		lErr.Fatalln(err.Error())
		os.Exit(1)
	}

	app.welcome()

	if app.Config.DataFile != "" {
		if err := app.loadDataFile(app.Config.DataFile); err != nil {
			// fmt.Fprintln(os.Stderr, err)
			lErr.Println(err.Error())
		}
	}

	if app.Config.BatchMode {
		if err := app.checkConfig(); err != nil {
			// fmt.Fprintln(os.Stderr, err.Error())
			lErr.Println(err.Error())
			flag.PrintDefaults()
			os.Exit(1)
		}

		if !app.isDataLoaded() {
			// fmt.Fprintln(os.Stderr, errors.New("File not loaded"))
			lErr.Println(errors.New("File not loaded"))
			os.Exit(1)
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("[%s]%s> ", app.prompt, func(b bool) string {
				if b {
					return "*"
				} else {
					return ""
				}
			}(app.isDataLoaded()))
			cmdString, err := reader.ReadString('\n')
			if err != nil {
				// fmt.Fprintln(os.Stderr, err)
				lErr.Fatalln(err.Error())
			}
			if err = app.runCommand(cmdString); err != nil {
				// fmt.Fprintln(os.Stderr, err)
				lErr.Println(err.Error())
			}
		}
	}

}
