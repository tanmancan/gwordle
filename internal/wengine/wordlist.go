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
	Definitions map[string]DictionaryApiDefinition
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

// Returns a cached definition for the given word. If no cache found, fetches and caches the definition first.
func (wl *WordList) GetDefinition(word string) (DictionaryApiDefinition) {
	if def, cached := wl.Definitions[word]; cached {
		return def
	}

	apiResponse := getWordDefinition(word)

	if (apiResponse.Error != DictionaryApiResponseError{}) {
		fmt.Println(apiResponse.Error)
		return nil
	}

	wl.Definitions[word] = apiResponse.Response[0]
	return wl.Definitions[word]
}

// Uses dictionaryapi.dev to see if the provided word is valid.
func (wl *WordList) CheckDictionary(word string) bool {
	definition := wl.GetDefinition(word)

	if (definition == DictionaryApiDefinition{}) {
		return false
	}

	return true
}

// Output the definition to the console
func (wl *WordList) ShowDefinition(word string) bool {
	definition := wl.GetDefinition(word)
	fmt.Printf(
		"%s (%s): %s\n",
		definition.Word,
		definition.Meanings[0].PartOfSpeech,
		definition.Meanings[0].Definitions[0].Definition,
	)
}
