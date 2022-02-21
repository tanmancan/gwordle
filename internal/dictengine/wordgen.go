package dictengine

import (
	"bufio"
	"fmt"
	"math/rand"
	"time"

	"github.com/tanmancan/gwordle/v1/internal/fhandler"
)

// Create a list of words grouped by their length.
type WordList struct {
	words map[int][]string // The key value is the length of the words in the value.
}

// Get a random secret word matching the requested length.
func GetSecretWord(length int) string {
	wordList := loadWordList()
	words := wordList.words[length]
	wordCount := len(words)
	rand.Seed(time.Now().UnixMilli())
	randomIdx := rand.Intn(wordCount)
	word := words[randomIdx]
	fmt.Println(wordCount, randomIdx, word)
	return word
}


// Parses a scanner generated from a word list file and returns a list of words
func ParseWordList(scanner *bufio.Scanner) (wordList WordList) {
	scanner.Split(bufio.ScanLines)
	wordList.words = make(map[int][]string)
	for scanner.Scan() {
		word := scanner.Text()
		length := len(word)
		wordList.words[length] = append(wordList.words[length], word)
	}

	return wordList
}

// Load the wordlist seeder and parse the words in groups based on word length.
func loadWordList() (wordList WordList) {
	wordListPath := "internal/dictengine/words_alpha.txt"
	reader := fhandler.WordListFileReader(wordListPath)
	scanner := fhandler.WordListFileScanner(reader)

	scanner.Split(bufio.ScanLines)
	wordList = ParseWordList(scanner)

	return wordList
}
