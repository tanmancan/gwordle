package gengine

import (
	"os"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"github.com/tanmancan/gwordle/v1/internal/dictionaryapi"
	"github.com/tanmancan/gwordle/v1/internal/localization"
	"github.com/tanmancan/gwordle/v1/internal/wengine"
)

type UserPrompt interface {
	// Get user input for a guess word or a help command.
	GetUserInput(gs *GameState) string
	// Displays help text.
	DisplayHelpText(gs *GameState)
	// Display a message when user loses a round.
	LoseRoundMessage(gs *GameState)
	// Display a message when a user wins a round.
	WinRoundMessage(gs *GameState)
	// Display a message when a user exists the game.
	ExitGameMessage(gs *GameState)
}

type Renderer interface {
	// Renders the result of the word validation for the current round.
	RenderValidationResults(gs *GameState)
	// Renders the current game score.
	RenderGameScore(gs *GameState)
	// Renders text inline,with string formatting.
	RenderText(format string, replacements ...interface{})
	// Renders text and adds a new line to the add, with string formatting.
	RenderTextLn(format string, replacements ...interface{})
}

type GameRound struct {
	Attempts int // Number of attemp remaining. Initial value is determined by AppConfig.UserConfig.MaxTries.
	Results []wengine.ValidationResult // Validation result for each guess word.
	SecretWord string // Current secret word.
	Win bool // If the current round was won.
}

type GameState struct {
	CurrentGame GameRound
	PastGames []GameRound
	GuessedWords []string
	UserPrompt UserPrompt
	Renderer Renderer
}

// Get the total number of wins and losses
func (gs *GameState) GetTotalWinLossCount() (win int, loss int) {
	for _, round := range gs.PastGames {
		if round.Win == true {
			win++
		} else {
			loss++
		}
	}

	return win, loss
}

// Initialize and begins the game loop
func (gs *GameState) InitGame(up UserPrompt, r Renderer) {
	gs.UserPrompt = up
	gs.Renderer = r
	gs.NewRound()
	gs.GameLoop()
}

// Main game loop that validates a users input and decides the outcome
func (gs *GameState) GameLoop() {
	completed := false
	for !completed {
		gs.Renderer.RenderValidationResults(gs)
		if gs.CurrentGame.Attempts == 0 {
			completed = true
			gs.LoseRound()
		}
		guess := gs.UserPrompt.GetUserInput(gs)
		completed = gs.ValidateGuessWord(guess)
	}
	gs.Renderer.RenderValidationResults(gs)
	gs.WinRound()
}

// Exit the game
func (gs *GameState) ExitGame() {
	// TODO: save and exit
	gs.UserPrompt.ExitGameMessage(gs)
	os.Exit(0)
}

// Validates the guess word
func (gs *GameState) ValidateGuessWord(word string) bool {
	if !wengine.WordListCache.HasWord(word) {
		request := dictionaryapi.GetWordDefinitionRequest{
			Word: word,
		}
		response := dictionaryapi.GetWordDefinition(request)
		if response.Response != nil && response.Response[0].Word == word {
			missingPath := "internal/cli/static/missing"
			wengine.WordListFileWriter(missingPath, word)
		} else {
			gs.Renderer.RenderTextLn(localization.AppTranslatable.Validation.InvalidWord, word)
			return false
		}
	}

	result, err := wengine.ValidateWord(word, gs.CurrentGame.SecretWord)

	if err != nil {
		gs.Renderer.RenderTextLn("%v", err)
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

// Start a new round with a new guess word
func (gs *GameState) NewRound() {
	gs.CurrentGame.SecretWord = wengine.WordListCache.GetRandomWord(config.GlobalConfig.UserConfig.WordLength)
	gs.CurrentGame.Attempts = config.GlobalConfig.UserConfig.MaxTries
	gs.CurrentGame.Results = nil
	gs.CurrentGame.Win = false
}

// Set win condition for the current round
func (gs *GameState) WinRound() {
	gs.UserPrompt.WinRoundMessage(gs)
	wengine.WordListCache.SetFilterWord(gs.CurrentGame.SecretWord)
	gs.CurrentGame.Win = true
	gs.PastGames = append(gs.PastGames, gs.CurrentGame)
	gs.Renderer.RenderGameScore(gs)
	gs.NewRound()
}

// Set lose condition for the current round.
func (gs *GameState) LoseRound() {
	gs.UserPrompt.LoseRoundMessage(gs)
	wengine.WordListCache.SetFilterWord(gs.CurrentGame.SecretWord)
	gs.PastGames = append(gs.PastGames, gs.CurrentGame)
	gs.Renderer.RenderGameScore(gs)
	gs.NewRound()
}
