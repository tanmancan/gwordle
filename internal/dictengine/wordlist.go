package dictengine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"github.com/tanmancan/gwordle/v1/internal/fhandler"
)

// Create a list of words grouped by their length.
type WordList struct {
	Words map[int][]string // The key value is the length of the words in the value.
}

// Response from dictionaryapi.dev
type DictionaryApiResponse struct {
	Word string
	Phonetic string
	Meanings []struct {
		PartOfSpeech string
		Definitions []struct {
			Definition string
		}
	}
}

// An error response from dictionaryapi.dev
type DictionaryApiResponseError struct {
	Title string
	Message string
	Resolution string
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
		invalidPath := fmt.Sprintf("internal/wloader/static/%s/invalid", config.GlobalConfig.Language.String())
		fhandler.WordListFileWriter(invalidPath, word)
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
	client := http.Client{}
	endpoint := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)
	request, reqErr := http.NewRequest("GET", endpoint, nil)

	if reqErr != nil {
		fmt.Printf("Error: %v\n", reqErr)
		return false
	}

	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false
	}

	content, errRead := ioutil.ReadAll(response.Body)

	if errRead != nil {
		fmt.Printf("Error: %v\n", errRead)
		return false
	}

	response.Body.Close()

	if response.StatusCode != 200 {
		var definitionError DictionaryApiResponseError
		json.Unmarshal(content, &definitionError)
		return false
	}

	var definition []DictionaryApiResponse
	json.Unmarshal(content, &definition)

	fmt.Printf("%s (%s): %s\n", definition[0].Word, definition[0].Meanings[0].PartOfSpeech, definition[0].Meanings[0].Definitions[0].Definition)

	return true
}
