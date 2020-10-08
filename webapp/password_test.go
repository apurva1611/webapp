package main

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestIsPasswordValid(t *testing.T) {
	tests := []struct {
		name  string
		pass  string
		valid bool
	}{
		{
			"NoCharacterAtAll",
			"",
			false,
		},
		{
			"JustEmptyStringAndWhitespace",
			" \n\t\r\v\f ",
			false,
		},
		{
			"MixtureOfEmptyStringAndWhitespace",
			"U u\n1\t?\r1\v2\f34",
			false,
		},
		{
			"MissingUpperCaseString",
			"uu1?1234",
			false,
		},
		{
			"MissingLowerCaseString",
			"UU1?1234",
			false,
		},
		{
			"MissingNumber",
			"Uua?aaaa",
			false,
		},
		{
			"MissingSymbol",
			"Uu101234",
			false,
		},
		{
			"LessThanRequiredMinimumLength",
			"Uu1?123",
			false,
		},
		{
			"ValidPassword",
			"Uu1?1234",
			true,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			if c.valid != IsPasswordValid(c.pass) {
				t.Fatal("invalid password")
			}
		})
	}
}
