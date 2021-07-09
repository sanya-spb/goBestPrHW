package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	lErr *log.Logger
	lOut *log.Logger
)

func main() {
	// Init our app
	app, err := newApp()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if fAccess, err := os.OpenFile(app.Config.LogAccess, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.Fatalln(err.Error())
	} else {
		defer fAccess.Close()
		lOut = log.New(fAccess, "", log.LstdFlags)
		lOut.Println("run")
	}
	if fErrors, err := os.OpenFile(app.Config.LogErrors, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.Fatalln(err.Error())
	} else {
		defer fErrors.Close()
		lErr = log.New(fErrors, "", log.LstdFlags)
		lErr.Println("run")
	}

	app.welcome()

	if app.Config.DataFile != "" {
		if err := app.loadDataFile(app.Config.DataFile); err != nil {
			lErr.Println(err.Error())
		}
	}

	if app.Config.BatchMode {
		if err := app.checkConfig(); err != nil {
			lErr.Println(err.Error())
			flag.PrintDefaults()
			os.Exit(1)
		}

		if !app.isDataLoaded() {
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
				lErr.Fatalln(err.Error())
			}
			if err = app.runCommand(cmdString); err != nil {
				fmt.Fprintln(os.Stderr, err)
				lErr.Println(err.Error())
			}
		}
	}
}

func string2Interface(value string, base int, bitSize int) interface{} {
	var result interface{}
	if i64, err := strconv.ParseInt(value, base, bitSize); err == nil {
		result = i64
	} else if f64, err := strconv.ParseFloat(value, bitSize); err == nil {
		result = f64
	} else {
		if len(value) > 0 && value[0] == '"' {
			value = value[1:]
		}
		if len(value) > 0 && value[len(value)-1] == '"' {
			value = value[:len(value)-1]
		}
		result = value
	}
	return result
}
