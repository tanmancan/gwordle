package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tanmancan/gwordle/v1/internal/dictengine"
)

// Maximum number of tries
var MaxTries int

// Worth length
var WordLength int

func main() {
	var results []dictengine.ValidationResult
	secret := dictengine.GetSecretWord(WordLength)
	userInput(secret, MaxTries, &results)
}

func init() {
	flag.IntVar(&MaxTries, "tries", 6, "Maximum number of tries. Default is 6.")
	flag.IntVar(&WordLength, "wlen", 5, "The word length. Default is 5")
	flag.Parse()
}

// Main game loop for the CLI application.
func userInput(secret string, tries int, results *[]dictengine.ValidationResult) {
	fmt.Printf("You have %d tries\n", tries)

	var guess string
	_, err := fmt.Scanln(&guess)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result, errResult := dictengine.ValidateWord(guess, secret)

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
		totalTries := MaxTries - tries
		if (totalTries == 1) {
			triesLabel = "try"
		}
		fmt.Printf("You have guessed the correct word: %s, in %v %s!\n", secret, totalTries, triesLabel)
		os.Exit(0)
	}

}

// Output the results of the guess word.
func displayValidation(result dictengine.ValidationResult) {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"

	for _, c := range result.Chars {
		var color string
		char := c.Char
		status := c.Status

		switch status {
		case dictengine.InvalidCharacter:
			color = colorReset
		case dictengine.InvalidPosition:
			color = colorYellow
		case dictengine.ValidPosition:
			color = colorGreen
		}

		fmt.Print(string(color), char, " ", string(colorReset))
	}

	fmt.Print("\n")
}
