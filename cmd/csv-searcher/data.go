package main

import (
	"errors"
	"fmt"
	"strings"
)

// structure for storing data
type Data struct {
	Head []string
	Data map[string]interface{}
	rows int
}

func (data *Data) setHead(headers []string) {
	for _, value := range headers {
		data.Head = append(data.Head, strings.TrimSpace(value))
	}
	data.rows = 0
	data.Data = nil
	// data.Head = headers
}

func (data *Data) getHead() []string {
	return data.Head
}

func (data *Data) addRow(row []interface{}) error {
	if len(row) != len(data.Head) {
		return errors.New("Columns in row not equal header")
	}

	for i, v := range data.Head {
		// t := data.Data[v]
		row[i] = strings.TrimSpace(fmt.Sprintf("%v", row[i]))
		switch t := data.Data[v].(type) {
		case []interface{}:
			data.Data[v] = append(t, row[i])
		case nil:
			col := []interface{}{row[i]}
			data.Data[v] = col
		}
	}
	data.rows++
	return nil
}

func (data *Data) selectHead(cols []string) {
	// fmt.Printf("DEBUG1: %+v\n", cols)
	for _, col := range cols {
		// fmt.Printf("DEBUG2: %s\n", col)
		for _, valH := range data.Head {
			// fmt.Printf("DEBUG3: %s\n", valH)
			if col == valH {
				fmt.Printf("%v\t", valH)
			}
		}
	}
	fmt.Printf("\n")
}

func (data *Data) selectRow(cols []string, row int) {
	for _, col := range cols {
		for _, valH := range data.Head {
			if col == valH {
				fmt.Printf("%v\t", data.Data[valH].([]interface{})[row])
			}
		}
		// fmt.Printf("%v\t", col)
	}
	fmt.Printf("\n")
}

func (data *Data) selectData(cols []string) {
	data.selectHead(cols)

	for ii := 0; ii < data.rows; ii++ {
		data.selectRow(cols, ii)
	}
	// for _, hVal := range data.Head {
	// 	for _, dVal := range data.Data[hVal] {
	// 		// fmt.Printf("%v\t", value)
	// 	}
	// 	fmt.Printf("\n")
	// }
	// return data.Head.([]interface{})
}
