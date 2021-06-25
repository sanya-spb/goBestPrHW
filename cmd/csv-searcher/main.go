package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sanya-spb/goBestPrHW/pkg_tools/config"
	"github.com/sanya-spb/goBestPrHW/pkg_tools/version"
)

type App struct {
	Version  version.AppVersion
	Config   config.Config
	DataFile string
	// logrus  *logrus.Logger
}

// Checking the required condition
func (app *App) checkConfig() error {
	if app.Config.DataFile == "" {
		return fmt.Errorf("No set path to *.cvs file")
	}
	return nil
}

// Checking the file with data
func (app *App) checkDataFile() error {
	return nil
}

func (app *App) loadDataFile(path string) error {
	if err := app.checkDataFile(); err != nil {
		return err
	}
	app.DataFile = path
	return nil
}

func (app *App) runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	switch arrCommandStr[0] {
	case "load":
		return app.loadDataFile(arrCommandStr[1])
	case "exit":
		os.Exit(0)
		// add another case here for custom commands.
	default:
		return fmt.Errorf("Unknown command!")
	}
	// cmd := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)
	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout
	// return cmd.Run()
	return nil
}

func newApp() *App {
	var app *App = new(App)
	app.Version = *version.Version
	app.Config = *config.NewConfig()
	return app
}

func main() {
	// Init our app
	app := newApp()

	fmt.Fprintf(os.Stderr, "%v\n", app)

	if app.Config.BatchMode {
		if err := app.checkConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			flag.PrintDefaults()
			os.Exit(1)
		}

		if err := app.checkDataFile(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("[csv-searcher]> ")
			cmdString, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			if err = app.runCommand(cmdString); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}
