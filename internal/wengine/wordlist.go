package wengine

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"github.com/tanmancan/gwordle/v1/internal/dictionaryapi"
)

// Create a list of words grouped by their length.
type WordList struct {
	Words map[int][]string // The key value is the length of the words in the value.
	Definitions map[string]dictionaryapi.DictionaryApiDefinition // Stores the definition for a word using api.dictionaryapi.dev
	FilterWords []string
}

var WordListCache WordList

// Get a random word from the wordlist that matches the request word length.
// Will filter out any words found within WordList.FilterWords
func (wl *WordList) GetRandomWord(length int) string {
	if len(wl.Words) == 0 {
		log.Fatalln("No words were loaded")
	}
	wordsPreFilter := wl.Words[length]

	if wordsPreFilter == nil {
		log.Fatalln("No word found for given length:", length)
	}

	words := wl.FilterWordList(wordsPreFilter)
	wl.Words[length] = words
	wordCount := len(words)

	rand.Seed(time.Now().UnixNano())
	randomIdx := 0
	if (wordCount > 1) {
		randomIdx = rand.Intn(wordCount - 1)
	}
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

// Add a word to the filter list
func (wl *WordList) SetFilterWord(word string) {
	existingIdx := sort.SearchStrings(wl.FilterWords, word)
	if !(existingIdx < len(wl.FilterWords) && wl.FilterWords[existingIdx] == word) {
		wl.FilterWords = append(wl.FilterWords, word)
		sort.Strings(wl.FilterWords)
	}
}

// Check if given word is in the filter list.
func (wl *WordList) HasFilterWord(word string) bool {
	if len(wl.FilterWords) == 0 {
		return false
	}

	sort.Strings(wl.FilterWords)
	searchIdx := sort.SearchStrings(wl.FilterWords, word)

	return searchIdx < len(wl.FilterWords) && wl.FilterWords[searchIdx] == word
}

// Apply WordList.FilterWords to the given list of words.
// Returns a sorted list with the filter words removed.
func (wl *WordList) FilterWordList(words []string) []string {
	if len(words) == 0 {
		return words
	}

	var filteredList []string

	for _, word := range words {
		if !wl.HasFilterWord(word) {
			filteredList = append(filteredList, word)
		}
	}

	sort.Strings(filteredList)

	return filteredList
}

// Returns a cached definition for the given word. If no cache found, fetches and caches the definition first.
func (wl *WordList) GetDefinition(word string) (*dictionaryapi.DictionaryApiDefinition) {
	if def, cached := wl.Definitions[word]; cached {
		return &def
	}

	request := dictionaryapi.GetWordDefinitionRequest{
		Word: word,
	}
	apiResponse := dictionaryapi.GetWordDefinition(request)

	if (apiResponse.Error != dictionaryapi.DictionaryApiResponseError{}) {
		return nil
	}

	if wl.Definitions == nil {
		wl.Definitions = make(map[string]dictionaryapi.DictionaryApiDefinition)
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
