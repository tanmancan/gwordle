package main

import (
	"fmt"
	"os"

	"github.com/tanmancan/gwordle/v1/internal/dictengine"
)

const maxTries int = 6

func main() {
	var results []dictengine.ValidationResult
	secret := dictengine.GetSecretWord(5)
	fmt.Println("Word is: ", secret)
	userInput(secret, maxTries, &results)
}

func userInput(secret string, tries int, results *[]dictengine.ValidationResult) {
	fmt.Printf("You have %d tries\n", tries)
	tries = tries - 1
	if (tries == 1) {
		fmt.Println("You loose")
		os.Exit(0)
	}

	var guess string
	_, err := fmt.Scanln(&guess)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result, errResult := dictengine.ValidateWord(guess, secret)
	*results = append(*results, result)

	for _, r := range *results {
		displayValidation(r)
	}

	if (result.Match == false) {
		userInput(secret, tries, results)
	} else {
		fmt.Printf("You have guessed the correct word: %s, in %v tries!\n", secret, maxTries - tries)
	}

	if errResult != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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
