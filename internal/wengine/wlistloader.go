package wengine

import (
	"bufio"
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"golang.org/x/text/language"
)

// Encapsulate an embedded file system where we store the word list files that are loaded into memory.
type WordFileSystem struct {
	fs embed.FS // The embed.FS filesystem where we load the original text file.
	validFilePathTemplate string // Template string to the list of valid words within the embed.FS filesystem. Example: static/%/valid
	invalidFilePathTemplate string // Template string to the list of invalid words within the embed.FS filesystem. Example: static/%/invalid
	locale language.Tag // The current language. By default this is used to generate the file paths. Example static/en/valid
}

// Get the constructed file path to the valid word list. Will try to apply the language to the template strings.
func (wfp *WordFileSystem) GetValidFilePath() string {
	return fmt.Sprintf(wfp.validFilePathTemplate, wfp.locale.String())
}

// Get the constructed fiel path to the invalid word list. Will try to apply the language to the template strings.
func (wfp *WordFileSystem) GetInvalidFilePath() string {
	return fmt.Sprintf(wfp.invalidFilePathTemplate, wfp.locale.String())
}

func (wfp *WordFileSystem) GetFileSystem() embed.FS {
	return wfp.fs
}

//go:embed static
var efs embed.FS

func init() {
	wordFs := WordFileSystem{
		fs: efs,
		validFilePathTemplate: "static/%s/valid",
		invalidFilePathTemplate: "static/%s/invalid",
		locale: config.GlobalConfig.Locale,
	}
	WordListCache = loadWordList(wordFs)
}

// Parses a scanner generated from a word list file and returns a list of words
func scanWordListFile(scannerValidList *bufio.Scanner, scannerInvalidList *bufio.Scanner) (wordList WordList) {
	var (
		validWordList []string
		invalidWordList []string
	)
	wordList.Words = make(map[int][]string)
	scannerValidList.Split(bufio.ScanLines)
	scannerInvalidList.Split(bufio.ScanLines)

	for scannerInvalidList.Scan() {
		wiv := strings.Trim(scannerInvalidList.Text(), " ")
		invalidWordList = append(invalidWordList, wiv)
	}

	sort.Strings(invalidWordList)

	for scannerValidList.Scan() {
		wv := strings.Trim(scannerValidList.Text(), " ")
		validWordList = append(validWordList, wv)
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
func loadWordList(wfs WordFileSystem) (wordList WordList) {
	validListReader, err := WordListFileReader(wfs.GetValidFilePath(), wfs.GetFileSystem())
	if (err != nil) {
		return wordList
	}

	invalidListReader, err := WordListFileReader(wfs.GetInvalidFilePath(), wfs.GetFileSystem())
	if (err != nil) {
		return wordList
	}

	scannerValidList := WordListFileScanner(validListReader)
	scannerInvalidList := WordListFileScanner(invalidListReader)

	return scanWordListFile(scannerValidList, scannerInvalidList)
}
