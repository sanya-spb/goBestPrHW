package main

import (
	"context"
	"sort"
	"testing"

	// "github.com/tj/assert"
	"github.com/stretchr/testify/assert"
)

func Test_Data_GetAllHeaders(t *testing.T) {
	type Test struct {
		data Data
		wait []string
	}

	ttEqual := []Test{
		{
			data: Data{
				Headers: Headers{},
				Data:    map[string]interface{}{},
				rows:    0,
			},
			wait: nil,
		},
		{
			data: Data{
				Headers: Headers{
					&Header{
						name:   "a",
						lenght: 1,
					},
				},
				Data: map[string]interface{}{},
				rows: 0,
			},
			wait: []string{
				"a",
			},
		},
	}

	for _, v := range ttEqual {
		assert.Equal(t, v.wait, v.data.getAllHeaders())
	}
}

func Test_Data_SetHead(t *testing.T) {
	type Test struct {
		data Data
		arg1 []string
		wait []string
	}

	ttEqual := []Test{
		{
			data: Data{
				Headers: Headers{},
				Data:    map[string]interface{}{},
				rows:    0,
			},
			arg1: []string{},
			wait: nil,
		},
		{
			data: Data{
				Headers: Headers{},
				Data:    map[string]interface{}{},
				rows:    0,
			},
			arg1: []string{
				"a",
				"b",
			},
			wait: []string{
				"a",
				"b",
			},
		},
	}

	for _, v := range ttEqual {
		v.data.setHead(v.arg1)
		assert.Equal(t, v.wait, v.data.getAllHeaders())
	}
}

func Test_Data_isHeader(t *testing.T) {
	type Test struct {
		data Data
		arg1 string
		wait bool
	}

	ttEqual := []Test{
		{
			data: Data{
				Headers: Headers{},
				Data:    map[string]interface{}{},
				rows:    0,
			},
			arg1: "",
			wait: false,
		},
		{
			data: Data{
				Headers: Headers{
					&Header{
						name:   "a",
						lenght: 1,
					},
				},
				Data: map[string]interface{}{},
				rows: 0,
			},
			arg1: "a",
			wait: true,
		},
	}

	for _, v := range ttEqual {
		assert.Equal(t, v.wait, v.data.isHeader(v.arg1))
	}
}

func Test_Data_addRow(t *testing.T) {
	type Test struct {
		data Data
		arg1 []interface{}
		wait Data
	}

	ttEqual := []Test{
		{
			data: Data{
				Headers: Headers{},
				Data:    map[string]interface{}{},
				rows:    0,
			},
			arg1: []interface{}{},
			wait: Data{
				Headers: Headers{},
				Data:    map[string]interface{}{},
				rows:    0,
			},
		},
		{
			data: Data{
				Headers: Headers{
					{
						name:   "int",
						lenght: 1,
					},
					{
						name:   "txt",
						lenght: 1,
					},
					{
						name:   "float",
						lenght: 1,
					},
				},
				Data: map[string]interface{}{},
				rows: 0,
			},
			arg1: []interface{}{
				int64(42),
				"abc",
				float64(3.1415),
			},
			wait: Data{
				Headers: Headers{
					{
						name:   "int",
						lenght: 2,
					},
					{
						name:   "txt",
						lenght: 3,
					},
					{
						name:   "float",
						lenght: 6,
					},
				},
				Data: map[string]interface{}{
					"int":   []interface{}{int64(42)},
					"txt":   []interface{}{"abc"},
					"float": []interface{}{float64(3.1415)},
				},
				rows: 1,
			},
		},
	}

	for _, v := range ttEqual {
		assert.NoError(t, v.data.addRow(v.arg1))
		assert.Equal(t, v.wait, v.data)
	}
}

func Test_Data_filterData(t *testing.T) {
	type Test struct {
		data   Data
		rows   []int
		filter Filter
		wait   []int
	}

	ttEqual := []Test{
		{
			data: Data{
				Headers: Headers{
					{
						name:   "a",
						lenght: 1,
					},
				},
				Data: map[string]interface{}{
					"a": []interface{}{int64(1)},
				},
				rows: 1,
			},
			rows: []int{0},
			filter: Filter{
				preposition: "",
				columnName:  "a",
				operator:    "=",
				value:       int64(1),
			},
			wait: []int{0},
		},
		{
			data: Data{
				Headers: Headers{
					{
						name:   "a",
						lenght: 1,
					},
					{
						name:   "b",
						lenght: 3,
					},
					{
						name:   "c",
						lenght: 3,
					},
				},
				Data: map[string]interface{}{
					"a": []interface{}{int64(1), int64(2), int64(3)},
					"b": []interface{}{float64(1.0), float64(2.0), float64(3.0)},
					"c": []interface{}{"1", "2", "3"},
				},
				rows: 3,
			},
			rows: []int{0, 1, 2},
			filter: Filter{
				preposition: "",
				columnName:  "a",
				operator:    "=",
				value:       int64(2),
			},
			wait: []int{1},
		},
	}

	ctx := context.Background()

	for _, v := range ttEqual {
		assert.NoError(t, v.data.filterData(ctx, &v.rows, v.filter))
		sort.Slice(v.rows, func(i, j int) bool {
			return v.rows[i] < v.rows[j]
		})
		assert.Equal(t, v.wait, v.rows)
	}
}

func Test_Data_runFilter(t *testing.T) {
	type Test struct {
		data    Data
		rows    []int
		filters []Filter
		wait    []int
	}

	ttEqual := []Test{
		{
			data: Data{
				Headers: Headers{
					{
						name:   "a",
						lenght: 1,
					},
					{
						name:   "b",
						lenght: 3,
					},
					{
						name:   "c",
						lenght: 3,
					},
				},
				Data: map[string]interface{}{
					"a": []interface{}{int64(1), int64(2), int64(3)},
					"b": []interface{}{float64(1.0), float64(2.0), float64(3.0)},
					"c": []interface{}{"1", "2", "3"},
				},
				rows: 3,
			},
			rows: []int{0, 1, 2},
			filters: []Filter{
				{
					preposition: "",
					columnName:  "a",
					operator:    "=",
					value:       int64(2),
				},
			},
			wait: []int{1},
		}, {
			data: Data{
				Headers: Headers{
					{
						name:   "a",
						lenght: 1,
					},
					{
						name:   "b",
						lenght: 3,
					},
					{
						name:   "c",
						lenght: 3,
					},
				},
				Data: map[string]interface{}{
					"a": []interface{}{int64(1), int64(2), int64(3)},
					"b": []interface{}{float64(1.0), float64(2.0), float64(3.0)},
					"c": []interface{}{"1", "2", "3"},
				},
				rows: 3,
			},
			rows: []int{0, 1, 2},
			filters: []Filter{
				{
					preposition: "",
					columnName:  "a",
					operator:    "=",
					value:       int64(2),
				},
				{
					preposition: "or",
					columnName:  "b",
					operator:    "=",
					value:       float64(1.0),
				},
			},
			wait: []int{0, 1},
		},
	}

	ctx := context.Background()

	for _, v := range ttEqual {
		assert.NoError(t, v.data.runFilter(ctx, &v.rows, v.filters))
		sort.Slice(v.rows, func(i, j int) bool {
			return v.rows[i] < v.rows[j]
		})
		assert.Equal(t, v.wait, v.rows)
	}
}
