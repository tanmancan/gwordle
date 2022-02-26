package config

import (
	"flag"

	"golang.org/x/text/language"
)

type AppConfig struct {
	Language language.Tag
	UserConfig UserConfig
}

type UserConfig struct {
	MaxTries int // The maximum number of guesses allowed in a game.
	WordLength int // The length of the guess word.
}

var GlobalConfig AppConfig

func init() {
	GlobalConfig.Language = language.English

	flag.IntVar(&GlobalConfig.UserConfig.MaxTries, "tries", 6, "Maximum number of tries. Default is 6.")
	flag.IntVar(&GlobalConfig.UserConfig.WordLength, "wlen", 5, "The word length. Default is 5")
}

func main() {
	flag.Parse()
}
