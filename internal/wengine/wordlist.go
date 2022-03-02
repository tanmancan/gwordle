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
	Definitions map[string]DictionaryApiDefinition // Stores the definition for a word using api.dictionaryapi.dev
	CompletedWords []string
}

var WordListCache WordList

// Get a random word from the wordlist that matches the request word length.
func (wl *WordList) GetRandomWord(length int) string {
	if len(wl.Words) == 0 {
		log.Fatalln("No words were loaded")
	}
	words := wl.Words[length]

	if words == nil {
		log.Fatalln("No word found for given length:", length)
	}

	wordCount := len(words)

	rand.Seed(time.Now().Unix())
	randomIdx := rand.Intn(wordCount - 1)
	word := words[randomIdx]

	valid := wl.CheckDictionary(word)

	if (!valid) {
		invalidPath := fmt.Sprintf("internal/wengine/static/%s/invalid", config.GlobalConfig.Locale.String())
		WordListFileWriter(invalidPath, word)
	} else {
		// wl.ShowDefinition(word)
	}

	return word
}

// Checks if the given word exists in the word list.
func (wl *WordList) HasWord(word string) bool {
	length := len(word)
	words := wl.Words[length]

	sort.Strings(words)
	searchIdx := sort.SearchStrings(words, word)

	return searchIdx < len(words) && words[searchIdx] == word
}

// func (wl *WordList) AddCompletedWord(word string) {
// 	existingIdx := sort.SearchStrings(wl.CompletedWords, word)
// }

// Returns a cached definition for the given word. If no cache found, fetches and caches the definition first.
func (wl *WordList) GetDefinition(word string) (*DictionaryApiDefinition) {
	if def, cached := wl.Definitions[word]; cached {
		return &def
	}

	request := buildDictionaryRequest(word)
	apiResponse := getWordDefinition(request)

	if (apiResponse.Error != DictionaryApiResponseError{}) {
		fmt.Println(apiResponse.Error)
		return nil
	}

	if wl.Definitions == nil {
		wl.Definitions = make(map[string]DictionaryApiDefinition)
	}
	def := apiResponse.Response[0]
	wl.Definitions[word] = def
	return &def
}

// Uses dictionaryapi.dev to see if the provided word is valid.
func (wl *WordList) CheckDictionary(word string) bool {
	definition := wl.GetDefinition(word)

	if (definition == nil) {
		return false
	}

	return true
}

// Output the definition to the console
func (wl *WordList) ShowDefinition(word string) {
	definition := wl.GetDefinition(word)

	if (definition == nil) {
		fmt.Println("No definition found for the word:", word)
	}
	fmt.Printf(
		"%s (%s): %s\n",
		definition.Word,
		definition.Meanings[0].PartOfSpeech,
		definition.Meanings[0].Definitions[0].Definition,
	)
}
