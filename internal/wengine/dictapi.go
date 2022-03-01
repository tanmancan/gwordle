package wengine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Response from api.dictionaryapi.dev
type DictionaryApiResponse struct {
	Response []DictionaryApiDefinition
	Error DictionaryApiResponseError
}

type DictionaryApiDefinition struct {
	Word string
	Phonetic string
	Phoenetics []struct {
		Text string
		Audio string
	}
	Origin string
	Meanings []struct {
		PartOfSpeech string
		Definitions []struct {
			Definition string
			Example string
		}
	}
}

// An error response from api.dictionaryapi.dev
type DictionaryApiResponseError struct {
	Title string
	Message string
	Resolution string
}

// Build a request for api.dictionaryapi.dev for the provided word.
func buildDictionaryRequest(word string) (*http.Request) {
	endpoint := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)
	request, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		log.Fatalln(err)
	}

	return request
}

// Parse the response from api.dictionaryapi.dev
func parseDictionaryResponse(response *http.Response) (apiResponse DictionaryApiResponse) {
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	switch response.StatusCode {
	case 200:
		json.Unmarshal(body, &apiResponse.Response)
	case 404:
		json.Unmarshal(body, &apiResponse.Error)
	}

	return apiResponse
}

// Get the definition for the provided word using api.dictionaryapi.dev
func getWordDefinition(word string) DictionaryApiResponse {
	request := buildDictionaryRequest(word)
	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	return parseDictionaryResponse(response)
}
