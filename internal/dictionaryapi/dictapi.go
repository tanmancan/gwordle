package dictionaryapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tanmancan/gwordle/v1/internal/config"
)

// Response from api.dictionaryapi.dev
type DictionaryApiResponse struct {
	Response []DictionaryApiDefinition
	Error DictionaryApiResponseError
}

type WordDefinitions struct {
	Definition string
	Example string
}

type WordMeanings struct {
	PartOfSpeech string
	Definitions []WordDefinitions
}

type DictionaryApiDefinition struct {
	Word string
	Phonetic string
	Origin string
	Meanings []WordMeanings
}

// An error response from api.dictionaryapi.dev
type DictionaryApiResponseError struct {
	Title string
	Message string
	Resolution string
}


// Build a request for api.dictionaryapi.dev for the provided word.
func BuildDictionaryRequest(word string) (*http.Request) {
	endpoint := fmt.Sprintf(config.GlobalConfig.DictionaryApiEndpoint, word)
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
	default:
		apiResponse.Error = DictionaryApiResponseError{
			Title: "Unknown error while fetching definition.",
			Message: fmt.Sprintf("Status: %s - %s", response.Status, string(body)),
			Resolution: "Check https://github.com/meetDeveloper/freeDictionaryAPI/issues for any service issues.",
		}

	}

	return apiResponse
}

// Get the definition for the provided word using api.dictionaryapi.dev
func GetWordDefinition(request *http.Request) DictionaryApiResponse {
	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	return parseDictionaryResponse(response)
}
