package wengine

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

//go:embed test-mocks/dictionaryapimocks/success-response.json
var mockSuccessResponse string

//go:embed test-mocks/dictionaryapimocks/error-response.json
var mockErrorResponse string

func Test_getWordDefinition(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		fmt.Fprintln(w, string(mockSuccessResponse))
	}))
	defer testServer.Close()
	request, err := http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name string
		args args
		want DictionaryApiResponse
	}{
		{
			name: "Call api with valid word: smile",
			args: args{
				request: request,
			},
			want: DictionaryApiResponse{
				Response: []DictionaryApiDefinition{
					{

						Word: "smile",
						Phonetic: "/ˈsmaɪ.əl/",
						Meanings: []WordMeanings{
							{
								PartOfSpeech: "noun",
								Definitions: []WordDefinitions{
									{
										Definition: "A facial expression comprised by flexing the muscles of both ends of one's mouth, often showing the front teeth, without vocalisation, and in humans is a common involuntary or voluntary expression of happiness, pleasure, amusement or anxiety.",
										Example: "He always puts a smile on my face.",
									},
									{
										Definition: "Favour; propitious regard.",
										Example: "the smile of the gods",
									},
								},
							},
							{
								PartOfSpeech: "verb",
								Definitions: []WordDefinitions {
									{
										Definition: "To have (a smile) on one's face.",
										Example: "I don't know what he's smiling about.",
									},
									{
										Definition: "To express by smiling.",
										Example: "to smile consent, or a welcome",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWordDefinition(tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWordDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getWordDefinition_errorResponse_statusNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, string(mockErrorResponse))
	}))
	defer testServer.Close()
	request, err := http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name string
		args args
		want DictionaryApiResponse
	}{
		{
			name: "Call api with invalid word",
			args: args{
				request: request,
			},
			want: DictionaryApiResponse{
				Error: DictionaryApiResponseError{
					Title:  "No Definitions Found",
					Message:  "Sorry pal, we couldn't find definitions for the word you were looking for.",
					Resolution: "You can try the search again at later time or head to the web instead.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWordDefinition(tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWordDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getWordDefinition_errorResponse_statusUnknown(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Oops!")
	}))
	defer testServer.Close()
	request, err := http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name string
		args args
		want DictionaryApiResponse
	}{
		{
			name: "Call api with invalid word",
			args: args{
				request: request,
			},
			want: DictionaryApiResponse{
				Error: DictionaryApiResponseError{
					Title: "Unknown error while fetching definition.",
					Message:  "Status: 400 Bad Request - Oops!\n",
					Resolution: "Check https://github.com/meetDeveloper/freeDictionaryAPI/issues for any service issues.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWordDefinition(tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWordDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}
