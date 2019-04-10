package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {



	var urlSet = []struct {
		urlToTest  string
		isValid bool
	}{
		{"https://i.imgur.com/mj9xOGP.mp4", true},
		{"https://i.imgur.com/nSyy4Bs.gifv", true},

		{"", false},
		{"nSyy4Bs.gifv", false},
		{"some string", false},
	}
	
	for _, tt := range urlSet {

		t.Run(tt.urlToTest, func(t *testing.T) {
			result := ifGifvURL(tt.urlToTest)
			assert.Equal(t, result, tt.isValid)
		})

	}
	
}