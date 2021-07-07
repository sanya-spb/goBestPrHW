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

// application struct
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

// cmd: pwd (show current working directory)
func (app *App) cmdPWD() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}

// cmd: ls (list files in current directory)
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

// cmd: cd (change dir)
func (app *App) cmdCD(dir string) error {
	if err := os.Chdir(dir); err != nil {
		return err
	}
	if _, err := os.Getwd(); err != nil {
		return err
	}
	return nil
}

// cmd: select (show data based on expression)
func (app *App) cmdSELECT(expr string) error {
	var (
		cols   []string
		sWhere string
	)

	splitExpr := strings.Split(expr, "where")
	if len(splitExpr) == 1 || len(splitExpr) == 2 {
		sCols := splitExpr[0]
		cols = strings.Split(sCols, ",")
		for k, v := range cols {
			cols[k] = strings.TrimSpace(v)
			if cols[k] == "*" {
				cols = app.Data.getAllHeaders()
				break
			}
		}
		if len(splitExpr) == 2 {
			sWhere = splitExpr[1]
		}
	} else {
		return errors.New("Invalid expession!")
	}

	fmt.Printf("cols: %v\nwhere: %s\n", cols, sWhere)

	app.Data.selectAllData(cols)
	return nil
}

// load csv data to memory
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
		if err := app.Data.addRow(cvsRowI); err != nil {
			lErr.Printf(err.Error())
		}
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

// user iteraction
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
		return app.Data.cmdHeaders()
	case "dump":
		fmt.Printf("%+v\n", app.Data)
	case "select":
		// app.Data.selectData(arrCommandStr[1:])
		return app.cmdSELECT(strings.Join(arrCommandStr[1:], " "))
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

// init for App
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

// print welcome message
func (app *App) welcome() {
	fmt.Fprintf(os.Stdout, "Welcome to csv-searcher!\nWorking directory: %s\nVersion: %s [%s@%s]\nCopyright: %s\n\n", app.exPath, app.Version.Version, app.Version.Commit, app.Version.BuildTime, app.Version.Copyright)
}

// get flag of loaded data
func (app *App) isDataLoaded() bool {
	if app.DataFile != "" {
		return true
	} else {
		return false
	}
}
