package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "embed"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"github.com/tanmancan/gwordle/v1/internal/gengine"
	"github.com/tanmancan/gwordle/v1/internal/localization"
	"github.com/tanmancan/gwordle/v1/internal/wengine"
)

type CliMemoryCard struct {}

type CliUserPrompt struct {}

type CliRenderer struct {}

// Get the save filepath in the user's home directory.
// Savefile are versioned. Old saves may not work with newer versions.
func (mc CliMemoryCard) GetSaveFilePath() (string, error) {
	hdir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	sdir := fmt.Sprintf("%s/gwordle/", hdir)

	if err := os.MkdirAll(sdir, os.ModePerm); err != nil {
		return "", err
	}
	return fmt.Sprintf("%ssave-%s.json", sdir, config.GlobalConfig.Version), nil
}

// Load game from file.
func (mc CliMemoryCard) LoadGame() *gengine.SaveState {
	sf, err := mc.GetSaveFilePath()
	if err != nil {
		log.Println(err)
		return nil
	}
	data, err := os.ReadFile(sf)
	if err != nil {
		log.Println(err)
		return nil
	}
	s := &gengine.SaveState{}
	json.Unmarshal(data, &s)
	return s
}

// Save game to file.
func (mc CliMemoryCard) SaveGame(s *gengine.SaveState) {
	sf, err := mc.GetSaveFilePath()
	if err != nil {
		log.Println(err)
		return
	}
	f, err := os.OpenFile(sf, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	b, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
		return
	}
	f.Write(b)
}

//go:embed static/hide
var hideBlock string

// Hides the current game prompt with some dummy logs.
func (up CliUserPrompt) HideGame(gs *gengine.GameState) {
	gs.Renderer.RenderTextLn(hideBlock)
	gs.Renderer.RenderTextLn(localization.AppTranslatable.HideRound.Instructions)
	up.HidePrompt(gs)
}

// Checks user prompt to cancel or continue hide display.
func (up CliUserPrompt) HidePrompt(gs *gengine.GameState) string {
	hideRound := localization.AppTranslatable.HideRound
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		gs.Renderer.RenderTextLn(hideRound.InvalidInput)
		return up.HidePrompt(gs)
	}

	switch strings.ToLower(input) {
	case strings.ToLower(hideRound.Return):
		return ""
	case strings.ToLower(hideRound.Exit):
		gs.Renderer.RenderTextLn(hideRound.Exit)
		os.Exit(0)
	default:
		gs.Renderer.RenderTextLn("%s %s", hideRound.InvalidInput, hideRound.Instructions)
		up.HidePrompt(gs)
	}

	return input
}

	// Get user input for a guess word or a help command.
func (up CliUserPrompt) GetUserInput(gs *gengine.GameState) string {
	gs.Renderer.RenderTextLn(localization.AppTranslatable.UserPrompt.Instructions, localization.AppTranslatable.Commands.Help)
	gs.Renderer.RenderText(localization.AppTranslatable.UserPrompt.RemainingAttempts, gs.SaveState.CurrentGame.RemainingAttempts);
	var guess string
	_, err := fmt.Scan(&guess)

	if err != nil {
		log.Fatalln(err)
	}

	guess = strings.ToLower(guess)

	if guess[0:1] == "/" {
		up.ParseUserCommand(guess, gs)
		return gs.UserPrompt.GetUserInput(gs)
	}

	return guess
}

func (up CliUserPrompt) ParseUserCommand(ucmd string, gs *gengine.GameState) {
	cmds := localization.AppTranslatable.Commands
	switch strings.Trim(ucmd, "/") {
	case cmds.Score:
		gs.Renderer.RenderGameScore(gs)
	case cmds.New:
		gs.LoseRound()
	case cmds.Help:
		gs.UserPrompt.DisplayHelpText(gs)
	case cmds.Exit:
		gs.ExitGame()
	case cmds.Hide:
		up.HideGame(gs)
	default:
		gs.Renderer.RenderTextLn(cmds.InvalidCommand, ucmd)
		gs.UserPrompt.DisplayHelpText(gs)
	}
}

// Displays help text.
func (up CliUserPrompt) DisplayHelpText(gs *gengine.GameState) {
	cmds := localization.AppTranslatable.Commands
	gs.Renderer.RenderTextLn("\n%s", cmds.HelpTextIntro)
	gs.Renderer.RenderTextLn("/%s		%s", cmds.Help, cmds.HelpDesc)
	gs.Renderer.RenderTextLn("/%s		%s", cmds.Score, cmds.ScoreDesc)
	gs.Renderer.RenderTextLn("/%s		%s", cmds.New, cmds.NewDesc)
	gs.Renderer.RenderTextLn("/%s		%s", cmds.Hide, cmds.HideDesc)
	gs.Renderer.RenderTextLn("/%s		%s\n", cmds.Exit, cmds.ExitDesc)
}

// Display a message when user loses a round.
func (up CliUserPrompt) LoseRoundMessage(gs *gengine.GameState) {
	labelsEndRound := localization.AppTranslatable.EndRound
	lMsg := fmt.Sprintf(labelsEndRound.LoseMessage, strings.ToUpper(gs.SaveState.CurrentGame.SecretWord))
	msg := fmt.Sprintf("| %s |", lMsg)
	dWidth := len(msg)
	hRuleSlice := make([]string, dWidth)
	for i := range hRuleSlice {
		hRuleSlice[i] = "-"
	}
	hRule := strings.Join(hRuleSlice, "")
	gs.Renderer.RenderTextLn("\n%s\n%s\n%s", hRule, msg, hRule)
}

// Display a message when a user wins a round.
func (up CliUserPrompt) WinRoundMessage(gs *gengine.GameState) {
	labelsEndRound := localization.AppTranslatable.EndRound
	triesLabel := labelsEndRound.Tries
	totalTries := config.GlobalConfig.UserConfig.MaxTries - gs.SaveState.CurrentGame.RemainingAttempts
	if (totalTries == 1) {
		triesLabel = labelsEndRound.Try
	}
	gs.Renderer.RenderTextLn(labelsEndRound.WinMessage, gs.SaveState.CurrentGame.SecretWord, totalTries, triesLabel)
}

// Display a message when a user exists the game.
func (up CliUserPrompt) ExitGameMessage(gs *gengine.GameState) {
	gs.Renderer.RenderTextLn("Good bye!")
}


// Renders the result of the word validation for the current round.
func (r CliRenderer) RenderValidationResults(gs *gengine.GameState) {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"
	fmt.Print("\n")
	for _, result := range gs.SaveState.CurrentGame.Results {
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
	for i := 0; i < gs.SaveState.CurrentGame.RemainingAttempts; i++ {
		for i := 0; i < config.GlobalConfig.UserConfig.WordLength; i++ {
			fmt.Print("_ ")
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}
// Renders the current game score.
func (r CliRenderer) RenderGameScore(gs *gengine.GameState) {
	win, loss := gs.GetTotalWinLossCount()
	scrCard := localization.AppTranslatable.ScoreCard
	gs.Renderer.RenderText("\n")
	gs.Renderer.RenderTextLn(scrCard.TotalWin, win)
	gs.Renderer.RenderTextLn(scrCard.TotalLoss, loss)
	gs.Renderer.RenderText("\n")
}
// Renders text inline,with string formatting.
func (r CliRenderer) RenderText(format string, replacements ...interface{}) {
	fmt.Printf(format, replacements...)
}
// Renders text and adds a new line to the add, with string formatting.
func (r CliRenderer) RenderTextLn(format string, replacements ...interface{}) {
	f := fmt.Sprintln(format)
	fmt.Printf(f, replacements...)
}

func InitCliGame() {
	up := CliUserPrompt{}
	r := CliRenderer{}
	mc := CliMemoryCard{}
	game := gengine.GameState{}
	game.InitGame(up, r, mc)
}
