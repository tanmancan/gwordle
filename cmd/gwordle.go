package main

import (
	"fmt"
	"os"

	"github.com/tanmancan/gwordle/v1/internal/dictengine"
)

func main() {
	var tries int = 5
	var results []dictengine.ValidationResult
	fmt.Printf("You have %d tries\n", tries)
	userInput(tries, &results)
}

func userInput(tries int, results *[]dictengine.ValidationResult) {
	if (tries == 0) {
		fmt.Println("You loose")
		os.Exit(0)
	}

	secret := "swill"

	var guess string
	_, err := fmt.Scanln(&guess)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	result, errResult := dictengine.ValidateWord(guess, secret)
	*results = append(*results, result)

	fmt.Printf("You have %d tries\n", tries)
	for _, r := range *results {
		displayValidation(r)
	}

	if (result.Match == false) {
		tries = tries - 1
		userInput(tries, results)
	} else {
		fmt.Printf("You have guessed the correct word, %s, in %v tries!\n", secret, tries)
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
