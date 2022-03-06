package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"github.com/tanmancan/gwordle/v1/internal/dictionaryapi"
	"github.com/tanmancan/gwordle/v1/internal/localization"
	"github.com/tanmancan/gwordle/v1/internal/wengine"
)

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
func (gs *GameState) InitGame() {
	gs.NewRound()
	gs.GameLoop()
}

// Main game loop that validates a users input and decides the outcome
func (gs *GameState) GameLoop() {
	completed := false
	for !completed {
		gs.RenderResults()
		if gs.CurrentGame.Attempts == 1 {
			completed = true
			gs.LoseRound()
		}
		guess := gs.GetUserInput()
		completed = gs.ValidateGuessWord(guess)
	}
	gs.WinRound()
}

// Exit the game
func (gs *GameState) ExitGame() {
	fmt.Println("TODO: save and exit")
	os.Exit(0)
}

// Get the user input
func (gs *GameState) GetUserInput() string {
	renderTextLn(localization.AppTranslatable.UserPrompt.Instructions, localization.AppTranslatable.Commands.Help)
	renderText(localization.AppTranslatable.UserPrompt.RemainingAttempts, gs.CurrentGame.Attempts);
	var guess string
	_, err := fmt.Scan(&guess)

	if err != nil {
		log.Fatalln(err)
	}

	if guess[0:1] == "/" {
		gs.ParseUserCommand(guess)
		return gs.GetUserInput()
	}

	return guess
}

func (gs *GameState) ParseUserCommand(ucmd string) {
	cmds := localization.AppTranslatable.Commands
	switch strings.Trim(ucmd, "/") {
	case cmds.Score:
		gs.RenderGameScore()
	case cmds.New:
		gs.LoseRound()
	case cmds.Help:
		gs.ShowHelp()
	case cmds.Exit:
		gs.ExitGame()
	default:
		renderTextLn(cmds.InvalidCommand, ucmd)
		gs.ShowHelp()
	}
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
		}
		renderTextLn(localization.AppTranslatable.Validation.InvalidWord, word)
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

// Start a new round with a new guess word
func (gs *GameState) NewRound() {
	gs.CurrentGame.SecretWord = wengine.WordListCache.GetRandomWord(config.GlobalConfig.UserConfig.WordLength)
	gs.CurrentGame.Attempts = config.GlobalConfig.UserConfig.MaxTries
	gs.CurrentGame.Results = nil
	gs.CurrentGame.Win = false
}

// Set win condition for the current round
func (gs *GameState) WinRound() {
	labelsEndRound := localization.AppTranslatable.EndRound
	triesLabel := labelsEndRound.Tries
	totalTries := config.GlobalConfig.UserConfig.MaxTries - gs.CurrentGame.Attempts
	if (totalTries == 1) {
		triesLabel = labelsEndRound.Try
	}
	renderTextLn(labelsEndRound.WinMessage, gs.CurrentGame.SecretWord, gs.CurrentGame.Attempts, triesLabel)
	wengine.WordListCache.SetFilterWord(gs.CurrentGame.SecretWord)
	gs.CurrentGame.Win = true
	gs.PastGames = append(gs.PastGames, gs.CurrentGame)
	gs.RenderGameScore()
	gs.NewRound()
}

// Set lose condition for the current round.
func (gs *GameState) LoseRound() {
	labelsEndRound := localization.AppTranslatable.EndRound
	lMsg := fmt.Sprintf(labelsEndRound.LoseMessage, strings.ToUpper(gs.CurrentGame.SecretWord))
	msg := fmt.Sprintf("| %s |", lMsg)
	dWidth := len(msg)
	hRuleSlice := make([]string, dWidth)
	for i := range hRuleSlice {
		hRuleSlice[i] = "-"
	}
	hRule := strings.Join(hRuleSlice, "")
	renderTextLn("\n%s\n%s\n%s", hRule, msg, hRule)

	wengine.WordListCache.SetFilterWord(gs.CurrentGame.SecretWord)
	gs.PastGames = append(gs.PastGames, gs.CurrentGame)
	gs.RenderGameScore()
	gs.NewRound()
}

// Renders results of the word validations for the current round.
func (gs *GameState) RenderResults() {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"
	fmt.Print("\n")
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

			fmt.Print(string(color), strings.ToUpper(char), " ", string(colorReset))
		}

		fmt.Print("\n")
	}
	fmt.Print("\n")
}

// Display game score
func (gs *GameState) RenderGameScore() {
	win, loss := gs.GetTotalWinLossCount()
	scrCard := localization.AppTranslatable.ScoreCard
	renderText("\n")
	renderTextLn(scrCard.TotalWin, win)
	renderTextLn(scrCard.TotalLoss, loss)
	renderText("\n")
}

// Show the available commands
func (gs *GameState) ShowHelp() {
	cmds := localization.AppTranslatable.Commands
	renderTextLn("\n%s", cmds.HelpTextIntro)
	renderTextLn("/%s		%s", cmds.Help, cmds.HelpDesc)
	renderTextLn("/%s		%s", cmds.Score, cmds.ScoreDesc)
	renderTextLn("/%s		%s", cmds.New, cmds.NewDesc)
	renderTextLn("/%s		%s", cmds.Hide, cmds.HideDesc)
	renderTextLn("/%s		%s\n", cmds.Exit, cmds.ExitDesc)
}

// Render text inline, with optional string formatting. A wrapper for fmt.Printf
func renderText(format string, replacements ...interface{}) {
	fmt.Printf(format, replacements...)
}

// Same as renderText, but adds a new line to the end of the resulting output.
func renderTextLn(format string, replacements ...interface{}) {
	f := fmt.Sprintln(format)
	fmt.Printf(f, replacements...)
}
