package main

import (
	"testing"

	// "github.com/tj/assert"
	"github.com/stretchr/testify/assert"
)

func Test_string2Interface(t *testing.T) {
	type Test struct {
		arg1 string
		arg2 int
		arg3 int
		wait interface{}
	}

	ttEqual := []Test{
		{
			"1", 10, 64,
			int64(1),
		},
		{
			"3.1415", 10, 64,
			float64(3.1415),
		},
		{
			"abc", 10, 64,
			"abc",
		},
	}

	for _, v := range ttEqual {
		assert.Equal(t, v.wait, string2Interface(v.arg1, v.arg2, v.arg3))
	}
}
