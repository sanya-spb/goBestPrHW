package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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
		cols    []string
		filters []Filter
	)

	splitExpr := strings.Split(expr, "where")
	if len(splitExpr) == 1 || len(splitExpr) == 2 {
		sCols := splitExpr[0]
		tCols := strings.Split(sCols, ",")
		for _, v := range tCols {
			v = strings.TrimSpace(v)
			if app.Data.isHeader(v) {
				cols = append(cols, v)
			} else if v == "*" {
				cols = app.Data.getAllHeaders()
				break
			} else {
				return errors.New("invalid column name!")
			}

		}
		if len(splitExpr) == 2 {
			sWhere := splitExpr[1]
			re1 := regexp.MustCompile(`(\s*(and|or){0,1}\s*(\w+\s*[><=]\s*[\w\"\.]+)\s*){1}`)
			re2 := regexp.MustCompile(`\s*(and|or){0,1}\s*(\w+)\s*([><=])\s*([\w\"\.]+)\s*`)
			for _, val := range re1.FindAllString(sWhere, -1) {
				for _, val2 := range re2.FindAllString(val, -1) {
					ff := re2.FindAllStringSubmatch(val2, -1)
					if len(ff[0]) >= 4 {
						filter := Filter{
							preposition: ff[0][len(ff[0])-4:][0],
							columnName:  ff[0][len(ff[0])-4:][1],
							operator:    ff[0][len(ff[0])-4:][2],
							value:       string2Interface(ff[0][len(ff[0])-4:][3], 10, 64),
						}
						if app.Data.isHeader(filter.columnName) {
							filters = append(filters, filter)
						} else {
							return errors.New("invalid column name!")
						}
					} else {
						return errors.New("Invalid filter!")
					}
				}
			}
		}
	} else {
		return errors.New("Invalid expession!")
	}

	if len(filters) == 0 {
		app.Data.selectAllData(cols)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.Config.FilterTimeout)*time.Millisecond)
		defer cancel()
		if err := app.Data.selectData(ctx, cols, filters); err != nil {
			return err
		}
	}
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
		cvsRowI := make([]interface{}, len(cvsRow))
		for i, txt := range cvsRow {
			cvsRowI[i] = string2Interface(strings.TrimSpace(txt), 10, 64)
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
	lOut.Printf("CMD: %s", commandStr)
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
		return app.cmdSELECT(strings.Join(arrCommandStr[1:], " "))
	case "exit":
		os.Exit(0)
		// add another case here for custom commands.
	default:
		return fmt.Errorf("Unknown command!")
	}
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
