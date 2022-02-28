package wengine

import (
	"sort"
	"strings"
)

// Metadata about the individual characters in the guessword.
type CharMetadata struct {
	Char string
	CountInGuess int // Number of times the character appears in the guess word.
	CountInSecret int // Number of times the character appears in the secret word.
	IndexesInGuess []int // The indexes of the character in the guess word. Nil if not found.
	IndexesInSecret []int // The indexes of the character in secret word. Nil if not found.
	IndexesValidGuess []int // The indexes in the guess words that were valid guesses. Nil if not found.
	IndexesInvalidGuess []int // The index in the guess words that were invalid guesses. Nil if not found.
}

// Adds a new index to the IndexesInGuess list
func (cm *CharMetadata) SetIndexInGuess(index int) *CharMetadata {
	cm.IndexesInGuess = append(cm.IndexesInGuess, index)
	sort.Ints(cm.IndexesInGuess)
	return cm
}

// Adds a new index to the IndexesInSecret list
func (cm *CharMetadata) SetIndexInSecret(index int) *CharMetadata {
	cm.IndexesInSecret = append(cm.IndexesInSecret, index)
	sort.Ints(cm.IndexesInSecret)
	return cm
}

// Adds a new index to the IndexesValidGuess list
func (cm *CharMetadata) SetIndexValidGuess(index int) *CharMetadata {
	cm.IndexesValidGuess = append(cm.IndexesValidGuess, index)
	sort.Ints(cm.IndexesValidGuess)
	return cm
}

// Adds a new index to the IndexesInvalidGuess list
func (cm *CharMetadata) SetIndexInvalidGuess(index int) *CharMetadata {
	cm.IndexesInvalidGuess = append(cm.IndexesInvalidGuess, index)
	sort.Ints(cm.IndexesInvalidGuess)
	return cm
}

// Returns whether or not the current character exists in the secret word.
func (cm *CharMetadata) InSecretWord() bool {
	return cm.CountInSecret > 0
}

// Returns true if all occurrence of the current character has been guessed correctly
func (cm *CharMetadata) FoundAllSecretChar() bool {
	for _, idxSecret := range cm.IndexesInSecret {
		idxMatch := sort.SearchInts(cm.IndexesInGuess, idxSecret)
		if (!(idxMatch < len(cm.IndexesInGuess) && idxSecret == cm.IndexesInGuess[idxMatch])) {
			return false
		}
	}

	return true
}

// Metadata about the guess word.
type WordMetadata struct {
	Chars map[string]CharMetadata
}

// Generate useful metadata for each character in the guess word compared to the secret word
func (gwm *WordMetadata) GenerateWordMetadata(guess string, secret string) {
	gwm.Chars = make(map[string]CharMetadata)
	guessChars := strings.Split(guess, "")
	secretChars := strings.Split(secret, "")

	for idxCharGuess, charGuess := range guessChars {
		metadataValue, metadataKey := gwm.Chars[charGuess]
		if !metadataKey  {
			metadataValue := CharMetadata{
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
