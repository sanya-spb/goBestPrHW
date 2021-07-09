package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// type Datas interface {
// 	getAllHeaders() []string
// 	isHeader(test string) bool
// }

// description of column
type Header struct {
	name   string
	lenght int
	// colType string
}

// slice for all column
type Headers []*Header

// structure for storing data
type Data struct {
	Headers
	Data map[string]interface{}
	rows int
}

type Filter struct {
	preposition string
	columnName  string
	operator    string
	value       interface{}
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
}

// get all column names
func (data *Data) getAllHeaders() []string {
	var result []string
	for _, val := range data.Headers {
		result = append(result, val.name)
	}
	return result
}

// check if exist this header
func (data *Data) isHeader(test string) bool {
	for _, header := range data.getAllHeaders() {
		if test == header {
			return true
		}
	}
	return false
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

	if len(row) == 0 {
		return nil
	}

	for key, value := range data.Headers {
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

// print row
func (data *Data) selectRow(cols []string, row int) {
	for _, col := range cols {
		for _, valH := range data.Headers {
			if col == valH.name {
				fmt.Printf("%-"+fmt.Sprint(valH.lenght+1)+"v", data.Data[valH.name].([]interface{})[row])
			}
		}
	}
	fmt.Printf("\n")
}

// simple version to get all data without filter
func (data *Data) selectAllData(cols []string) {
	data.selectHead(cols)
	for ii := 0; ii < data.rows; ii++ {
		data.selectRow(cols, ii)
	}
}

// advanced version to get data with filter
func (data *Data) selectData(ctx context.Context, cols []string, filters []Filter) error {
	data.selectHead(cols)

	rows := make([]int, 0, data.rows)
	for ii := 0; ii < data.rows; ii++ {
		rows = append(rows, ii)
	}

	if err := data.runFilter(ctx, &rows, filters); err != nil {
		return err
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i] < rows[j]
	})

	for _, row := range rows {
		data.selectRow(cols, row)
	}
	return nil
}

// run filters step by step (with timeout)
func (data *Data) runFilter(ctx context.Context, rows *[]int, filters []Filter) error {
	select {
	case <-ctx.Done():
		return errors.New("Context is done!")
	default:
		if len(filters) > 0 {
			filter := filters[0]
			filters := filters[1:]
			if err := data.filterData(ctx, rows, filter); err != nil {
				return err
			}
			if err := data.runFilter(ctx, rows, filters); err != nil {
				return err
			}
		} else {
			return nil
		}
	}
	return nil
}

// filter data
func (data *Data) filterData(ctx context.Context, rows *[]int, filter Filter) error {
	if filter == (Filter{}) {
		return errors.New("Empty filter!")
	}
	if len(data.Headers) == 0 {
		return errors.New("Empty headers!")
	}
	if data.rows == 0 {
		return errors.New("Empty data!")
	}
	if *rows == nil {
		return errors.New("Empty input rows!")
	}
	index := make(map[int]bool, len(*rows))
	for _, v := range *rows {
		index[v] = true
	}
	switch filter.preposition {
	case "&&", "":
		for _, row := range *rows {
			select {
			case <-ctx.Done():
				return errors.New("Context is done!")
			default:
				switch t := data.Data[filter.columnName].([]interface{})[row].(type) {
				case string:
					if fmt.Sprintf("%T", filter.value) == "string" {
						switch filter.operator {
						case "=":
							if !(t == filter.value.(string)) {
								index[row] = false
							}
						case ">":
							if !(t > filter.value.(string)) {
								index[row] = false
							}
						case "<":
							if !(t < filter.value.(string)) {
								index[row] = false
							}
						default:
							return errors.New("undefined operator!")
						}
					} else {
						index[row] = false
					}
				case int64:
					if fmt.Sprintf("%T", filter.value) == "int64" {
						switch filter.operator {
						case "=":
							if !(t == filter.value.(int64)) {
								index[row] = false
							}
						case ">":
							if !(t > filter.value.(int64)) {
								index[row] = false
							}
						case "<":
							if !(t < filter.value.(int64)) {
								index[row] = false
							}
						default:
							return errors.New("undefined operator!")
						}
					} else {
						index[row] = false
					}
				case float64:
					if fmt.Sprintf("%T", filter.value) == "float64" {
						switch filter.operator {
						case "=":
							if !(t == filter.value.(float64)) {
								index[row] = false
							}
						case ">":
							if !(t > filter.value.(float64)) {
								index[row] = false
							}
						case "<":
							if !(t < filter.value.(float64)) {
								index[row] = false
							}
						default:
							return errors.New("undefined operator!")
						}
					} else {
						index[row] = false
					}
				}
			}
		}
	case "||":
		for row := 0; row < data.rows; row++ {
			select {
			case <-ctx.Done():
				return errors.New("Context is done!")
			default:
				switch t := data.Data[filter.columnName].([]interface{})[row].(type) {
				case string:
					if fmt.Sprintf("%T", filter.value) == "string" {
						switch filter.operator {
						case "=":
							if t == filter.value.(string) {
								index[row] = true
							}
						case ">":
							if t > filter.value.(string) {
								index[row] = true
							}
						case "<":
							if t < filter.value.(string) {
								index[row] = true
							}
						default:
							return errors.New("undefined operator!")
						}
					} else {
						index[row] = false
					}
				case int64:
					if fmt.Sprintf("%T", filter.value) == "int64" {
						switch filter.operator {
						case "=":
							if t == filter.value.(int64) {
								index[row] = true
							}
						case ">":
							if t > filter.value.(int64) {
								index[row] = true
							}
						case "<":
							if t < filter.value.(int64) {
								index[row] = true
							}
						default:
							return errors.New("undefined operator!")
						}
					} else {
						index[row] = false
					}
				case float64:
					if fmt.Sprintf("%T", filter.value) == "float64" {
						switch filter.operator {
						case "=":
							if t == filter.value.(float64) {
								index[row] = true
							}
						case ">":
							if t > filter.value.(float64) {
								index[row] = true
							}
						case "<":
							if t < filter.value.(float64) {
								index[row] = true
							}
						default:
							return errors.New("undefined operator!")
						}
					} else {
						index[row] = false
					}
				}
			}
		}
	}
	result := make([]int, 0, len(*rows))
	for k, v := range index {
		select {
		case <-ctx.Done():
			return errors.New("Context is done!")
		default:
			if v {
				result = append(result, k)
			}
		}
	}
	*rows = result
	return nil
}
