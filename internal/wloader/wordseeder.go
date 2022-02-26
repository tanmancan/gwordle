package wloader

import (
	"bufio"
	"embed"
	"fmt"
	"sort"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"github.com/tanmancan/gwordle/v1/internal/dictengine"
	"github.com/tanmancan/gwordle/v1/internal/fhandler"
)

// World list files are stored within the static directory, grouped by languages.
// For each languages, there are two files - valid and invalid.
// All words within valid will be loaded as possible game words.
// Any words that exists within the invalid file, will be fileted out while loading the valid wordlist.
//
// Each wordlist must be stored within their own language directory. Each directory must use
// the canonical string names for the corresponding language.Tag.
//
// Example: language.English.String() = "en"
//
// Example directory structure
//
// static/en/valid
//
// static/en/invalid
//go:embed static/*
var wordsFs embed.FS

func init() {
	dictengine.WordListCache = loadWordList(config.GlobalConfig.Language.String(), wordsFs)
}

// Parses a scanner generated from a word list file and returns a list of words
func scanWordListFile(scannerValidList *bufio.Scanner, scannerInvalidList *bufio.Scanner) (wordList dictengine.WordList) {
	var validWordList []string
	var invalidWordList []string
	wordList.Words = make(map[int][]string)
	scannerValidList.Split(bufio.ScanLines)
	scannerInvalidList.Split(bufio.ScanLines)

	for scannerInvalidList.Scan() {
		invalidWordList = append(invalidWordList, scannerInvalidList.Text())
	}

	sort.Strings(invalidWordList)

	for scannerValidList.Scan() {
		validWordList = append(validWordList, scannerValidList.Text())
	}

	sort.Strings(validWordList)

	for _, word := range validWordList {
		if len(invalidWordList) > 0 {
			invalidIdx := sort.SearchStrings(invalidWordList, word)
			if invalidIdx < len(invalidWordList) && word == invalidWordList[invalidIdx] {
				continue
			}
		}
		length := len(word)
		wordList.Words[length] = append(wordList.Words[length], word)
	}

	return wordList
}

// Load the wordlist seeder and parse the words in groups based on word length.
func loadWordList(language string, fs embed.FS) (wordList dictengine.WordList) {
	validListPath := fmt.Sprintf("static/%s/valid", language)
	invalidListPath := fmt.Sprintf("static/%s/invalid", language)

	validListReader, err := fhandler.WordListFileReader(validListPath, fs)
	if (err != nil) {
		return wordList
	}

	invalidListReader, err := fhandler.WordListFileReader(invalidListPath, fs)
	if (err != nil) {
		return wordList
	}

	scannerValidList := fhandler.WordListFileScanner(validListReader)
	scannerInvalidList := fhandler.WordListFileScanner(invalidListReader)

	return scanWordListFile(scannerValidList, scannerInvalidList)
}
