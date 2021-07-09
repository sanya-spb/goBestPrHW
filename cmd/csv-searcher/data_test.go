package main

import (
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
				42,
				"abc",
				3.1415,
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
					"int":   []interface{}{42},
					"txt":   []interface{}{"abc"},
					"float": []interface{}{3.1415},
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
