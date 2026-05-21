package main

import (
	"strings"
)

// split the input into "words" based on whitespace
// lowercase the input and trim all  whitespace
func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
