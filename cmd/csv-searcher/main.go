package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sanya-spb/goBestPrHW/internal/config"
	"github.com/sanya-spb/goBestPrHW/pkg/version"
)

type App struct {
	Version  version.AppVersion
	Config   config.Config
	DataFile string
	Data     map[string]interface{}
	DataHead []string
	exPath   string
	lErr     *log.Logger
	lOut     *log.Logger
	prompt   string
}

// Checking the required condition
func (app *App) checkConfig() error {
	if app.Config.DataFile == "" {
		return fmt.Errorf("No set path to *.cvs file")
	}
	return nil
}

// Checking the file with data
// func (app *App) checkDataFile() error {
// 	if _, err := os.Stat(app.DataFile); !os.IsNotExist(err) {
// 		return err
// 	}
// 	return nil
// }

func (app *App) loadDataFile(path string) error {
	if _, err := os.Stat(app.DataFile); !os.IsNotExist(err) {
		return err
	}

	fData, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fData.Close()

	scanner := bufio.NewScanner(fData)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	app.DataHead = strings.Split(scanner.Text(), ",")

	// var txtlines []string
	// []interface{}

	// for scanner.Scan() {
	// 	txtlines := strings.Split(scanner.Text(), ",")
	// 	for ii := 0; ii < len(app.DataHead); ii++ {
	// 		app.Data[app.DataHead[ii]] = append(app.Data[app.DataHead[ii]], txtlines[ii])
	// 	}
	// }

	app.DataFile = path
	return nil
}

// func (app *App) getExecPath() (exPath string, err error) {
// 	ex, err := os.Executable()
// 	if err != nil {
// 		return "", err
// 	}
// 	exPath = filepath.Dir(ex)
// 	return exPath, nil
// }

func (app *App) runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	if len(arrCommandStr) <= 0 {
		return nil
	}
	switch arrCommandStr[0] {
	case "load":
		return app.loadDataFile(arrCommandStr[1])
	case "config":
		fmt.Printf("%+v\n", app.Config)
	case "headers":
		fmt.Printf("%+v\n", app.DataHead)
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

func newApp() (*App, error) {
	var app *App = new(App)
	app.Version = *version.Version
	app.Config = *config.NewConfig()
	app.prompt = "csv-searcher"

	if ex, err := os.Executable(); err != nil {
		return nil, err
	} else {
		ex = filepath.Dir(ex)
		app.exPath = ex
	}

	if fAccess, err := os.OpenFile(app.Config.LogAccess, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		return nil, errors.New(err.Error())
	} else {
		defer fAccess.Close()
		app.lOut = log.New(fAccess, "", log.LstdFlags)
		app.lOut.Println("run")
	}
	if fErrors, err := os.OpenFile(app.Config.LogErrors, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		return nil, errors.New(err.Error())
	} else {
		defer fErrors.Close()
		app.lErr = log.New(fErrors, "", log.LstdFlags)
		app.lErr.Println("run")
	}

	return app, nil
}

func (app *App) welcome() {
	fmt.Fprintf(os.Stderr, "Welcome to csv-searcher!\nWorking directory: %s\nVersion: %s [%s@%s]\nCopyright: %s\n\n", app.exPath, app.Version.Version, app.Version.Commit, app.Version.BuildTime, app.Version.Copyright)
}

func (app *App) isDataLoaded() bool {
	if app.DataFile != "" {
		return true
	} else {
		return false
	}
}

func main() {
	// Init our app
	app, err := newApp()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	app.welcome()

	if app.Config.DataFile != "" {
		if err := app.loadDataFile(app.Config.DataFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	if app.Config.BatchMode {
		if err := app.checkConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			flag.PrintDefaults()
			os.Exit(1)
		}

		if !app.isDataLoaded() {
			fmt.Fprintln(os.Stderr, errors.New("File not loaded"))
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
				fmt.Fprintln(os.Stderr, err)
			}
			if err = app.runCommand(cmdString); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}

}
