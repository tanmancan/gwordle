package config

import (
	"flag"

	"golang.org/x/text/language"
)

type appConfig struct {
	// Current application language. Default is English.
	Locale language.Tag
	// User configurations.
	UserConfig userConfig
	// API endpoint for dictionary lookup.
	DictionaryApiEndpoint string
	// Current application version.
	// TODO: Inject at build time automatically.
	Version string
}

type userConfig struct {
	MaxTries int // The maximum number of guesses allowed in a game.
	WordLength int // The length of the guess word.
}

var GlobalConfig appConfig

func init() {
	GlobalConfig.Version = "0.0.1"
	GlobalConfig.Locale = language.English
	GlobalConfig.DictionaryApiEndpoint = "https://api.dictionaryapi.dev/api/v2/entries/en/%s"
	flag.IntVar(&GlobalConfig.UserConfig.MaxTries, "tries", 6, "Maximum number of tries. Default is 6.")
	flag.IntVar(&GlobalConfig.UserConfig.WordLength, "wlen", 5, "The word length. Default is 5")
}

func main() {
	flag.Parse()
}
