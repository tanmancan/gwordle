package dictengine

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
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

// Load the wordlist seeder and parse the words in groups based on word length.
func loadWordList() (wordList WordList) {
	data, err := os.ReadFile("internal/dictengine/words_alpha.txt")
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	scanner.Split(bufio.ScanLines)
	wordList.words = make(map[int][]string)

	for scanner.Scan() {
		word := scanner.Text()
		length := len(word)
		wordList.words[length] = append(wordList.words[length], word)
	}

	if (err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}

	return wordList
}
