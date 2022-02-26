package fhandler

import (
	"bufio"
	"embed"
	"strings"
)

// Get a strings.Reader for the provided text file path
func WordListFileReader(path string, fs embed.FS) (*strings.Reader, error) {
	data, err := fs.ReadFile(path)

	if (err != nil) {
		return nil, err
	}

	return strings.NewReader(string(data)), nil
}

// Generate a Scanner from the provided strings.Reader
func WordListFileScanner(reader *strings.Reader) (*bufio.Scanner) {
	var scanner = &bufio.Scanner{}
	scanner = bufio.NewScanner(reader)
	return scanner
}
