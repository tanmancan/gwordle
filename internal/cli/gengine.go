package cli

import (
	"fmt"
	"log"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"github.com/tanmancan/gwordle/v1/internal/dictionaryapi"
	"github.com/tanmancan/gwordle/v1/internal/wengine"
)

// type GameEngine interface {
// 	InitGame() // Begins the game loop
// 	ExitGame() // Exits the game
// 	GetUserInput() string // Returns the guess word from the user
// 	ValidateUserInput(word string) bool // Validate the user guess
// 	RenderResults() // Display current results
// 	NewRound() // Starts a new game
// 	WinRound() // Correct word has been guessed
// 	LoseRound() // Ran out of attempts
// }

type GameRound struct {
	Attempts int // Number of attemp remaining. Initial value is determined by AppConfig.UserConfig.MaxTries.
	Results []wengine.ValidationResult
	SecretWord string
}

type GameState struct {
	CurrentGame GameRound
	PastGames []GameRound
	GuessedWords []string
}

func (gs *GameState) InitGame() {
	gs.NewRound()
	gs.GameLoop()
}

func (gs *GameState) GameLoop() {
	completed := false
	for !completed {
		gs.RenderResults()
		if gs.CurrentGame.Attempts == 1 {
			completed = true
			gs.LoseRound()
		}
		guess := gs.GetUserInput()
		completed = gs.ValidateUserInput(guess)
	}
	gs.WinRound()
}

func (gs *GameState) ExitGame() {}

func (gs *GameState) GetUserInput() string {
	fmt.Printf("You have %d tries\n", gs.CurrentGame.Attempts)
	var guess string
	_, err := fmt.Scanln(&guess)

	if err != nil {
		log.Fatalln(err)
	}

	return guess
}

func (gs *GameState) ValidateUserInput(word string) bool {
	if !wengine.WordListCache.HasWord(word) {
		request := dictionaryapi.GetWordDefinitionRequest{
			Word: word,
		}
		response := dictionaryapi.GetWordDefinition(request)
		if response.Response != nil && response.Response[0].Word == word {
			missingPath := "internal/cli/static/missing"
			wengine.WordListFileWriter(missingPath, word)
		}
		fmt.Println("Invalid word:", word)
		return false
	}

	result, err := wengine.ValidateWord(word, gs.CurrentGame.SecretWord)

	if err != nil {
		fmt.Println(err)
		return false
	}

	gs.CurrentGame.Attempts -= 1
	gs.CurrentGame.Results = append(gs.CurrentGame.Results, result)

	if result.Match == false {
		return false
	} else {
		return true
	}
}

func (gs *GameState) RenderResults() {
	fmt.Println("AHOY")
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"
	for _, result := range gs.CurrentGame.Results {
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
}

func (gs *GameState) NewRound() {
	gs.CurrentGame.SecretWord = wengine.WordListCache.GetRandomWord(config.GlobalConfig.UserConfig.WordLength)
	gs.CurrentGame.Attempts = config.GlobalConfig.UserConfig.MaxTries
}

func (gs *GameState) WinRound() {
	triesLabel := "tries"
	totalTries := config.GlobalConfig.UserConfig.MaxTries - gs.CurrentGame.Attempts
	if (totalTries == 1) {
		triesLabel = "try"
	}
	fmt.Printf("You have guessed the correct word: %s, in %v %s!\n", gs.CurrentGame.SecretWord, gs.CurrentGame.Attempts, triesLabel)
	wengine.WordListCache.SetFilterWord(gs.CurrentGame.SecretWord)
	gs.PastGames = append(gs.PastGames, gs.CurrentGame)
	gs.NewRound()
}

func (gs *GameState) LoseRound() {
	fmt.Println("You lose")
	fmt.Println("Word is: ", gs.CurrentGame.SecretWord)
	wengine.WordListCache.SetFilterWord(gs.CurrentGame.SecretWord)
	gs.PastGames = append(gs.PastGames, gs.CurrentGame)
	gs.NewRound()
}


