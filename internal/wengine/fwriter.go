package wengine

import (
	"fmt"
	"os"
)

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
