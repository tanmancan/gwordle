package main

import (
	"fmt"
	"os"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"github.com/tanmancan/gwordle/v1/internal/wengine"
)

func main() {
	var results []wengine.ValidationResult
	secret := wengine.WordListCache.GetRandomWord(config.GlobalConfig.UserConfig.WordLength)
	userInput(secret, config.GlobalConfig.UserConfig.MaxTries, &results)
}

// Main game loop for the CLI application.
func userInput(secret string, tries int, results *[]wengine.ValidationResult) {
	fmt.Printf("You have %d tries\n", tries)

	var guess string
	_, err := fmt.Scanln(&guess)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result, errResult := wengine.ValidateWord(guess, secret)

	if errResult != nil {
		fmt.Println(errResult)
		userInput(secret, tries, results)
	}

	tries = tries - 1
	if (tries == 1) {
		fmt.Println("You loose")
		fmt.Println("Word is: ", secret)
		os.Exit(0)
	}

	*results = append(*results, result)

	for _, r := range *results {
		displayValidation(r)
	}

	if (result.Match == false) {
		userInput(secret, tries, results)
	} else {
		triesLabel := "tries"
		totalTries := config.GlobalConfig.UserConfig.MaxTries - tries
		if (totalTries == 1) {
			triesLabel = "try"
		}
		fmt.Printf("You have guessed the correct word: %s, in %v %s!\n", secret, totalTries, triesLabel)
		os.Exit(0)
	}

}

// Output the results of the guess word.
func displayValidation(result wengine.ValidationResult) {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"

	for _, c := range result.Chars {
		var color string
		char := c.Char
		status := c.Status

		switch status {
		case wengine.InvalidCharacter:
			color = colorReset
		case wengine.InvalidPosition:
			color = colorYellow
		case wengine.ValidPosition:
			color = colorGreen
		}

		fmt.Print(string(color), char, " ", string(colorReset))
	}

	fmt.Print("\n")
}
