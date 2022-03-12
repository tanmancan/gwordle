package wengine

import (
	"errors"
	"sort"
	"strings"
)

// Determines if a character is in a valid position, invalid position, or is not a valid guess.
type CharValidationStatus = string

const (
	ValidPosition CharValidationStatus = "VALID_POS" // The correct character and position
	InvalidPosition = "INVALID_POS" // Correct character but invalid position
	InvalidCharacter = "INVALID_CHAR" // Invalid character guessed
)

// Compare the individual characters of the guess word to the secret word.
type CharValidationResult struct {
	Char string
	Status CharValidationStatus
}

// Validation result when comparing a guess word to a secret word.
type ValidationResult struct {
	Match bool // Does the guess word match the secret word.
	Chars []CharValidationResult
}


// Compares and validates a guess word against the secret word.
func ValidateWord(guess string, secret string) (result ValidationResult, err error) {
	guess = strings.ToLower(guess)
	secret = strings.ToLower(secret)
	result.Match = strings.Compare(guess, secret) == 0

	if len(guess) != len(secret) {
		err = errors.New("The guess and secret words are not the same length.")
		return result, err
	}

	var guessWordMetadata WordMetadata
	guessWordMetadata.GenerateWordMetadata(guess, secret)
	guessChars := strings.Split(guess, "")

	for i, c := range guessChars {
		var compStatus CharValidationStatus
		cMetadata := guessWordMetadata.Chars[c]
		repeatingIndexInGuess := sort.SearchInts(cMetadata.IndexesInGuess, i)

		switch {
		case secret[i] == guess[i]:
			compStatus = ValidPosition
		case !cMetadata.InSecretWord():
			compStatus = InvalidCharacter
		case cMetadata.FoundAllSecretChar() && cMetadata.InSecretWord() && secret[i] != guess[i]:
			compStatus = InvalidCharacter
		case len(cMetadata.IndexesInSecret) > 0 &&
		repeatingIndexInGuess < cMetadata.CountInSecret:
			compStatus = InvalidPosition
		default:
			compStatus = InvalidCharacter
		}

		result.Chars = append(result.Chars, CharValidationResult{
			Char: cMetadata.Char,
			Status: compStatus,
		})
	}

	return result, err
}

