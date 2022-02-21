package fhandler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Get a strings.Reader for the provided text file path
func WordListFileReader(path string) *strings.Reader {
	data, err := os.ReadFile(path)

	if (err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}

	return strings.NewReader(string(data))
}

// Generate a Scanner from the provided strings.Reader
func WordListFileScanner(reader *strings.Reader) (*bufio.Scanner) {
	var scanner = &bufio.Scanner{}
	scanner = bufio.NewScanner(reader)
	return scanner
}
