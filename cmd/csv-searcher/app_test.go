package main

import (
	"testing"

	// "github.com/tj/assert"
	"github.com/stretchr/testify/assert"
)

func Test_App_isDataLoaded(t *testing.T) {
	type Test struct {
		app  App
		wait bool
	}

	ttEqual := []Test{
		{
			app: App{
				DataFile: "",
			},
			wait: false,
		},
		{
			app: App{
				DataFile: "abc",
			},
			wait: true,
		},
	}

	for _, v := range ttEqual {
		assert.Equal(t, v.wait, v.app.isDataLoaded())
	}
}
