package localization

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tanmancan/gwordle/v1/internal/config"
	"golang.org/x/text/language"
)

type Translatables struct {
	// Slash commands. ex: /help /score
	Commands struct {
		Help string
		HelpDesc string
		Score string
		ScoreDesc string
		New string
		NewDesc string
		Hide string
		HideDesc string
		Exit string
		ExitDesc string
		InvalidCommand string
		HelpTextIntro string
	}
	UserPrompt struct {
		Instructions string
		RemainingAttempts string
	}
	ScoreCard struct {
		TotalWin string
		TotalLoss string
	}
	Validation struct {
		InvalidWord string
	}
	EndRound struct {
		Try string
		Tries string
		WinMessage string
		LoseMessage string
	}
	HideRound struct {
		Return string
		Exit string
		Instructions string
		InvalidInput string
	}
}

var (
	//go:embed static
	translatableFs embed.FS
	AppTranslatable Translatables
)

// Merge current language with base english text. This will ensure any missing
// labels are replaced with the original english values.
func init() {
	langEs := language.English.String()
	fPathEs := fmt.Sprintf("static/%s/translatables.json", langEs)
	translatableFileEs, err := translatableFs.ReadFile(fPathEs)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(translatableFileEs, &AppTranslatable)

	if config.GlobalConfig.Locale != language.English {
		lang := config.GlobalConfig.Locale.String()
		fPath := fmt.Sprintf("static/%s/translatables.json", lang)
		translatableFile, err := translatableFs.ReadFile(fPath)
		if err != nil {
			log.Fatalln(err)
		}
		json.Unmarshal(translatableFile, &AppTranslatable)
	}
}
