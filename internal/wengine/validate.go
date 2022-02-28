package wengine

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
	CountInSecret int // Number of times the character appears in the secret word.
	IndexesInGuess []int // The indexes of the character in the guess word. Nil if not found.
	IndexesInSecret []int // The indexes of the character in secret word. Nil if not found.
	IndexesValidGuess []int // The indexes in the guess words that were valid guesses. Nil if not found.
	IndexesInvalidGuess []int // The index in the guess words that were invalid guesses. Nil if not found.
}

// Adds a new index to the IndexesInGuess list
func (cm *GuessWordCharMetadata) SetIndexInGuess(index int) *GuessWordCharMetadata {
	cm.IndexesInGuess = append(cm.IndexesInGuess, index)
	sort.Ints(cm.IndexesInGuess)
	return cm
}

// Adds a new index to the IndexesInSecret list
func (cm *GuessWordCharMetadata) SetIndexInSecret(index int) *GuessWordCharMetadata {
	cm.IndexesInSecret = append(cm.IndexesInSecret, index)
	sort.Ints(cm.IndexesInSecret)
	return cm
}

// Adds a new index to the IndexesValidGuess list
func (cm *GuessWordCharMetadata) SetIndexValidGuess(index int) *GuessWordCharMetadata {
	cm.IndexesValidGuess = append(cm.IndexesValidGuess, index)
	sort.Ints(cm.IndexesValidGuess)
	return cm
}

// Adds a new index to the IndexesInvalidGuess list
func (cm *GuessWordCharMetadata) SetIndexInvalidGuess(index int) *GuessWordCharMetadata {
	cm.IndexesInvalidGuess = append(cm.IndexesInvalidGuess, index)
	sort.Ints(cm.IndexesInvalidGuess)
	return cm
}

// Returns whether or not the current character exists in the secret word.
func (cm *GuessWordCharMetadata) InSecretWord() bool {
	return cm.CountInSecret > 0
}

// Returns true if all occurrence of the current character has been guessed correctly
func (cm *GuessWordCharMetadata) FoundAllSecretChar() bool {
	for _, idxSecret := range cm.IndexesInSecret {
		idxMatch := sort.SearchInts(cm.IndexesInGuess, idxSecret)
		if (!(idxMatch < len(cm.IndexesInGuess) && idxSecret == cm.IndexesInGuess[idxMatch])) {
			return false
		}
	}

	return true
}

// Metadata about the guess word.
type GuessWordMetadata struct {
	Chars map[string]GuessWordCharMetadata
}

// Generate useful metadata for each character in the guess word compared to the secret word
func (gwm *GuessWordMetadata) GenerateGuessWordMetadata(guess string, secret string) {
	gwm.Chars = make(map[string]GuessWordCharMetadata)
	guessChars := strings.Split(guess, "")
	secretChars := strings.Split(secret, "")

	for idxCharGuess, charGuess := range guessChars {
		metadataValue, metadataKey := gwm.Chars[charGuess]
		if !metadataKey  {
			metadataValue := GuessWordCharMetadata{
				Char: charGuess,
				CountInGuess: strings.Count(guess, charGuess),
				CountInSecret: strings.Count(secret, charGuess),
			}

			metadataValue.SetIndexInGuess(idxCharGuess)

			for idxCharSecret, charSecret := range secretChars {
				if (charGuess == charSecret) {
					metadataValue.SetIndexInSecret(idxCharSecret)
					if (idxCharGuess == idxCharSecret) {
						metadataValue.SetIndexValidGuess(idxCharGuess)
					} else {
						metadataValue.SetIndexInvalidGuess(idxCharGuess)
					}
				}
			}
			gwm.Chars[charGuess] = metadataValue
		} else {
			metadataValue.SetIndexInGuess(idxCharGuess)
			gwm.Chars[charGuess] = metadataValue
		}
	}
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

	var guessWordMetadata GuessWordMetadata
	guessWordMetadata.GenerateGuessWordMetadata(guess, secret)
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

		result.Chars = append(result.Chars, CharacterValidationResult{
			Char: cMetadata.Char,
			Status: compStatus,
		})
	}

	return result, err
}

