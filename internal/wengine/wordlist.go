package wengine

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/tanmancan/gwordle/v1/internal/config"
)

// Create a list of words grouped by their length.
type WordList struct {
	Words map[int][]string // The key value is the length of the words in the value.
}

var WordListCache WordList

// Get a random word from the wordlist that matches the request word length.
func (wl *WordList) GetRandomWord(length int) string {
	if len(wl.Words) == 0 {
		log.Fatalln("No words were loaded")
	}
	words := wl.Words[length]
	wordCount := len(words)

	rand.Seed(time.Now().UnixMilli())
	randomIdx := rand.Intn(wordCount)
	word := words[randomIdx]
	// @todo Remove after Debugging.
	fmt.Println(wordCount, randomIdx, word)
	valid := wl.CheckDictionary(word)

	if (!valid) {
		invalidPath := fmt.Sprintf("internal/wengine/static/%s/invalid", config.GlobalConfig.Locale.String())
		WordListFileWriter(invalidPath, word)
	}

	return word
}

// Checks if the given word exists in the word list.
func (wl *WordList) HasWord(word string) bool {
	length := len(word)
	words := wl.Words[length]

	searchIdx := sort.SearchStrings(words, word)

	return searchIdx < len(words) && words[searchIdx] == word
}

// Uses dictionaryapi.dev to see if the provided word is valid.
func (wl *WordList) CheckDictionary(word string) bool {
	apiResponse := getWordDefinition(word)
	if (apiResponse.Error != DictionaryApiResponseError{}) {
		fmt.Println(apiResponse.Error)
		return false
	}

	fmt.Printf(
		"%s (%s): %s\n",
		apiResponse.Response[0].Word,
		apiResponse.Response[0].Meanings[0].PartOfSpeech,
		apiResponse.Response[0].Meanings[0].Definitions[0].Definition,
	)

	return true
}
