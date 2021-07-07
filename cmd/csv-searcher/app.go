package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
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
	Data
	exPath string
	prompt string
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

func (app *App) cmdPWD() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}

func (app *App) cmdLS() error {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			fmt.Printf("%s/\n", f.Name())
		} else {
			fmt.Printf("%s\n", f.Name())
		}
	}
	return nil
}

func (app *App) cmdCD(dir string) error {
	if err := os.Chdir(dir); err != nil {
		return err
	}
	if _, err := os.Getwd(); err != nil {
		return err
	}
	return nil
}

func (app *App) loadDataFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
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
	app.Data.setHead(strings.Split(scanner.Text(), ","))

	app.Data.Data = make(map[string]interface{})
	for scanner.Scan() {
		cvsRow := strings.Split(scanner.Text(), ",")
		//converting a []string to a []interface{}
		cvsRowI := make([]interface{}, len(cvsRow))
		for i, v := range cvsRow {
			cvsRowI[i] = v
		}
		_ = app.Data.addRow(cvsRowI)
	}

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
	case "pwd":
		return app.cmdPWD()
	case "ls":
		return app.cmdLS()
	case "cd":
		return app.cmdCD(arrCommandStr[1])
	case "load":
		return app.loadDataFile(arrCommandStr[1])
	case "config":
		fmt.Printf("%+v\n", app.Config)
	case "headers":
		fmt.Printf("%+v\n", app.Data.getHead())
	case "dump":
		fmt.Printf("%+v\n", app.Data)
	case "select":
		app.Data.selectData(arrCommandStr[1:])
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
		lOut = log.New(fAccess, "", log.LstdFlags)
		lOut.Println("run")
	}
	if fErrors, err := os.OpenFile(app.Config.LogErrors, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		return nil, errors.New(err.Error())
	} else {
		defer fErrors.Close()
		lErr = log.New(fErrors, "", log.LstdFlags)
		lErr.Println("run")
	}

	return app, nil
}

func (app *App) welcome() {
	fmt.Fprintf(os.Stdout, "Welcome to csv-searcher!\nWorking directory: %s\nVersion: %s [%s@%s]\nCopyright: %s\n\n", app.exPath, app.Version.Version, app.Version.Commit, app.Version.BuildTime, app.Version.Copyright)
}

func (app *App) isDataLoaded() bool {
	if app.DataFile != "" {
		return true
	} else {
		return false
	}
}
