package wengine

import (
	"bufio"
	"embed"
	"fmt"
	"os"
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

// Appends the given word to the given file path.
// If file does not exists, it will create it.
// Returns the number of bytes or error.
func WordListFileWriter(path string, word string) (int, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return 0, err
	}

	n, err := f.WriteString(fmt.Sprintf("%s\n", word))

	if err != nil {
		f.Close()
		return n, err
	}

	f.Close()
	return n, err
}
