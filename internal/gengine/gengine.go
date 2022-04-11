package gengine

import (
	"os"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"github.com/tanmancan/gwordle/v1/internal/dictionaryapi"
	"github.com/tanmancan/gwordle/v1/internal/localization"
	"github.com/tanmancan/gwordle/v1/internal/wengine"
)

// Used for user interaction and dialog.
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

// Used to display messages regarding the game state.
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

// Allows for saving and loading games.
type MemoryCard interface {
	// Loads an existing game.
	LoadGame() *SaveState
	// Save the game.
	SaveGame(s *SaveState)
}

// A single game round.
type GameRound struct {
	RemainingAttempts int // Number of attemp remaining. Initial value is determined by AppConfig.UserConfig.MaxTries.
	Results []wengine.ValidationResult // Validation result for each guess word.
	SecretWord string // Current secret word.
	Win bool // If the current round was won.
}

// The game state.
type GameState struct {
	SaveState SaveState
	UserPrompt UserPrompt
	Renderer Renderer
	MemoryCard MemoryCard
}

// Gamestate that can be saved and loaded
type SaveState struct {
	// Current round
	CurrentGame GameRound
	// Past rounds
	PastGames []GameRound
}

// Get the total number of wins and losses
func (gs *GameState) GetTotalWinLossCount() (win int, loss int) {
	for _, round := range gs.SaveState.PastGames {
		if round.Win == true {
			win++
		} else {
			loss++
		}
	}

	return win, loss
}

// Initialize and begins the game loop
func (gs *GameState) InitGame(up UserPrompt, r Renderer, mc MemoryCard) {
	gs.UserPrompt = up
	gs.Renderer = r
	gs.MemoryCard = mc
	prev := mc.LoadGame()
	if prev != nil {
		gs.SaveState = *prev
	} else {
		gs.NewRound()
	}
	gs.GameLoop()
}

// Main game loop that validates a users input and decides the outcome
func (gs *GameState) GameLoop() {
	completed := false
	for !completed {
		gs.Renderer.RenderValidationResults(gs)
		if gs.SaveState.CurrentGame.RemainingAttempts == 0 {
			completed = true
			gs.LoseRound()
		}
		guess := gs.UserPrompt.GetUserInput(gs)
		completed = gs.ValidateGuessWord(guess)
	}
	gs.Renderer.RenderValidationResults(gs)
	gs.WinRound()
	gs.GameLoop()
}

// Exit the game
func (gs *GameState) ExitGame() {
	gs.MemoryCard.SaveGame(&gs.SaveState)
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

	result, err := wengine.ValidateWord(word, gs.SaveState.CurrentGame.SecretWord)

	if err != nil {
		gs.Renderer.RenderTextLn("%v", err)
		return false
	}

	gs.SaveState.CurrentGame.RemainingAttempts -= 1
	gs.SaveState.CurrentGame.Results = append(gs.SaveState.CurrentGame.Results, result)

	if result.Match == false {
		return false
	} else {
		return true
	}
}

// Start a new round with a new guess word
func (gs *GameState) NewRound() {
	gs.SaveState.CurrentGame.SecretWord = wengine.WordListCache.GetRandomWord(config.GlobalConfig.UserConfig.WordLength)
	gs.SaveState.CurrentGame.RemainingAttempts = config.GlobalConfig.UserConfig.MaxTries
	gs.SaveState.CurrentGame.Results = nil
	gs.SaveState.CurrentGame.Win = false
}

// Set win condition for the current round
func (gs *GameState) WinRound() {
	gs.UserPrompt.WinRoundMessage(gs)
	wengine.WordListCache.SetFilterWord(gs.SaveState.CurrentGame.SecretWord)
	gs.SaveState.CurrentGame.Win = true
	gs.SaveState.PastGames = append(gs.SaveState.PastGames, gs.SaveState.CurrentGame)
	gs.Renderer.RenderGameScore(gs)
	gs.NewRound()
}

// Set lose condition for the current round.
func (gs *GameState) LoseRound() {
	gs.UserPrompt.LoseRoundMessage(gs)
	wengine.WordListCache.SetFilterWord(gs.SaveState.CurrentGame.SecretWord)
	gs.SaveState.PastGames = append(gs.SaveState.PastGames, gs.SaveState.CurrentGame)
	gs.Renderer.RenderGameScore(gs)
	gs.NewRound()
}
