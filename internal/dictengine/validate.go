package dictengine

import (
	"errors"
	"sort"
	"strings"
)

// Determines if a character is in a valid position, invalid position, or is not a valid guess.
type CharValidationStatus = int64

const (
	ValidPosition CharValidationStatus = iota // The correct character and position
	InvalidPosition // Correct character but invalid position
	InvalidCharacter // Invalid character guessed
)

// Compare the individual characters of the guess word to the secret word.
type CharacterValidationResult struct {
	Char string
	Status CharValidationStatus
}

// Validation result when comparing a guess word to a secret word.
type ValidationResult struct {
	Match bool // Does the guess word match the secret word.
	Chars []CharacterValidationResult
}

// Metadata about the individual characters in the guessword.
type GuessWordCharMetadata struct {
	Char string
	CountInGuess int // Number of times the character appears in the guess word.
	IndexesInGuess []int // First instance of the character in the guess word. -1 if not found.
	CountInSecret int // Number of times the character appears in the secret word.
	IndexesInSecret []int // First instance of the character in secret word. -1 if not found.
}

// Metadata about the guess word.
type GuessWordMetadata struct {
	Chars map[string]GuessWordCharMetadata
}

// Generate useful metadata for each character in the guess word compared to the secret word
func generateGuessWordMetadata(guess string, secret string) (metadata GuessWordMetadata) {
	metadata.Chars = make(map[string]GuessWordCharMetadata)

	guessChars := strings.Split(guess, "")
	secretChars := strings.Split(secret, "")

	for i, c := range guessChars {
		var indexesInSecret []int
		for idxCharSecret, charSecret := range secretChars {
			if (c == charSecret) {
				indexesInSecret = append(indexesInSecret, idxCharSecret)
			}
		}
		indexesInGuess := append(metadata.Chars[c].IndexesInGuess, i)

		metadata.Chars[c] = GuessWordCharMetadata{
			Char: c,
			CountInGuess: strings.Count(guess, c),
			IndexesInGuess: indexesInGuess,
			CountInSecret: strings.Count(secret, c),
			IndexesInSecret: indexesInSecret,
		}
	}

	return metadata
}

// Compares and validates a guess word against the secret word.
func ValidateWord(guess string, secret string) (result ValidationResult, err error) {
	guess = strings.ToLower(guess)
	secret = strings.ToLower(secret)
	result.Match = strings.Compare(guess, secret) == 0

	if len(guess) != len(secret) {
		err = errors.New("The guess and secret words are not the same length")
		return result, err
	}

	guessWordMetadata := generateGuessWordMetadata(guess, secret)
	guessChars := strings.Split(guess, "")

	for i, c := range guessChars {
		var compStatus CharValidationStatus
		cMetadata := guessWordMetadata.Chars[c]
		repeatingIndexInGuess := sort.SearchInts(cMetadata.IndexesInGuess, i)

		switch {
		case cMetadata.CountInSecret == 0:
			compStatus = InvalidCharacter
		case secret[i] == guess[i]:
			compStatus = ValidPosition
		case len(cMetadata.IndexesInSecret) > 0 &&
		repeatingIndexInGuess < cMetadata.CountInSecret:
			compStatus = InvalidPosition
		default:
			compStatus = InvalidCharacter
		}

		result.Chars = append(result.Chars, CharacterValidationResult{
			Char: cMetadata.Char,
			Status: compStatus,
		})
	}

	return result, err
}

