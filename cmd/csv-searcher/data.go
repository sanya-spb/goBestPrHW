package main

import (
	"errors"
	"fmt"
	"strings"
)

// description of column
type Header struct {
	name   string
	lenght int
}

// slice for all column
type Headers []*Header

// structure for storing data
type Data struct {
	Headers
	Data map[string]interface{}
	rows int
}

// fill header in data struct
func (data *Data) setHead(headers []string) {
	data.Headers = make(Headers, 0, len(headers))
	data.Data = make(map[string]interface{})
	data.rows = 0
	for _, value := range headers {
		var tHead Header
		tHead.name = strings.TrimSpace(value)
		tHead.lenght = len(tHead.name)
		data.Headers = append(data.Headers, &tHead)
	}
	// data.Head = headers
}

// get all column names
func (data *Data) getAllHeaders() []string {
	var result []string
	for _, val := range data.Headers {
		result = append(result, val.name)
	}
	return result
}

// execute cmd: headers
func (data *Data) cmdHeaders() error {
	var (
		maxValWidth int
		maxValLen   int
	)

	for _, value := range data.Headers {
		if len(value.name) > maxValWidth {
			maxValWidth = len(value.name)
		}
		if len(fmt.Sprint(value.lenght)) > maxValLen {
			maxValLen = len(fmt.Sprint(value.lenght))
		}
	}
	// sort.Sort(sHead)
	for _, valH := range data.Headers {
		fmt.Printf("%-"+fmt.Sprint(maxValWidth+1)+"s length: %"+fmt.Sprint(maxValLen+1)+"d\n", valH.name, valH.lenght)
	}
	return nil
}

// fill row in data struct
func (data *Data) addRow(row []interface{}) error {
	if len(row) != len(data.Headers) {
		return errors.New("Columns in row not equal header")
	}

	for key, value := range data.Headers {
		// t := data.Data[v]
		row[key] = strings.TrimSpace(fmt.Sprintf("%v", row[key]))
		if value.lenght < len(fmt.Sprint(row[key])) {
			value.lenght = len(fmt.Sprint(row[key]))
			data.Headers[key] = value
		}
		switch t := data.Data[value.name].(type) {
		case []interface{}:
			data.Data[value.name] = append(t, row[key])
		case nil:
			col := []interface{}{row[key]}
			data.Data[value.name] = col
		}
	}
	data.rows++
	return nil
}

// show header
func (data *Data) selectHead(cols []string) {
	for _, col := range cols {
		for _, valH := range data.Headers {
			if col == valH.name {
				fmt.Printf("%-"+fmt.Sprint(valH.lenght+1)+"v", valH.name)
			}
		}
	}
	fmt.Printf("\n")
}

// simple version to get row without filter
func (data *Data) selectAllRow(cols []string, row int) {
	for _, col := range cols {
		for _, valH := range data.Headers {
			if col == valH.name {
				fmt.Printf("%-"+fmt.Sprint(valH.lenght+1)+"v", data.Data[valH.name].([]interface{})[row])
			}
		}
		// fmt.Printf("%v\t", col)
	}
	fmt.Printf("\n")
}

// simple version to get all data without filter
func (data *Data) selectAllData(cols []string) {
	data.selectHead(cols)
	for ii := 0; ii < data.rows; ii++ {
		data.selectAllRow(cols, ii)
	}
}
