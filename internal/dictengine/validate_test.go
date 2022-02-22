package dictengine

import (
	"reflect"
	"testing"
)

func TestValidateWord(t *testing.T) {
	type args struct {
		guess  string
		secret string
	}
	tests := []struct {
		name       string
		args       args
		wantResult ValidationResult
		wantErr    bool
	}{
		{
			name: "Guess word: lxxxx. Secret word: swill",
			args: args{
				guess:  "lxxxx",
				secret: "swill",
			},
			wantResult: ValidationResult{
				Match: false,
				Chars: []CharacterValidationResult{
					{
						Char:   "l",
						Status: InvalidPosition,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Guess word: xlxxx. Secret word: swill",
			args: args{
				guess:  "xlxxx",
				secret: "swill",
			},
			wantResult: ValidationResult{
				Match: false,
				Chars: []CharacterValidationResult{
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "l",
						Status: InvalidPosition,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Guess word: xxxlx. Secret word: swill",
			args: args{
				guess:  "xxxlx",
				secret: "swill",
			},
			wantResult: ValidationResult{
				Match: false,
				Chars: []CharacterValidationResult{
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "l",
						Status: ValidPosition,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Guess word: xxxxl. Secret word: swill",
			args: args{
				guess:  "xxxxl",
				secret: "swill",
			},
			wantResult: ValidationResult{
				Match: false,
				Chars: []CharacterValidationResult{
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "l",
						Status: ValidPosition,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Guess word: wxwxx. Secret word: swill",
			args: args{
				guess:  "wxwxx",
				secret: "swill",
			},
			wantResult: ValidationResult{
				Match: false,
				Chars: []CharacterValidationResult{
					{
						Char:   "w",
						Status: InvalidPosition,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "w",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
					{
						Char:   "x",
						Status: InvalidCharacter,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := ValidateWord(tt.args.guess, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ValidateWord() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_generateGuessWordMetadata(t *testing.T) {
	type args struct {
		guess  string
		secret string
	}
	wantChar := make(map[string]GuessWordCharMetadata)
	wantChar["w"] = GuessWordCharMetadata{
		Char: "w",
		CountInGuess: 1,
		IndexesInGuess: []int{0},
		CountInSecret: 1,
		IndexesInSecret: []int{1},
	}
	wantChar["x"] = GuessWordCharMetadata{
		Char: "x",
		CountInGuess: 4,
		IndexesInGuess: []int{1,2,3,4},
		CountInSecret: 0,
		IndexesInSecret: []int{},
	}
	tests := []struct {
		name         string
		args         args
		wantMetadata GuessWordMetadata
	}{
		{
			name: "Genrate metadata",
			args: args{
				guess: "wxxxx",
				secret: "swill",
			},
			wantMetadata: GuessWordMetadata{
				Chars: wantChar,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMetadata := generateGuessWordMetadata(tt.args.guess, tt.args.secret); !reflect.DeepEqual(gotMetadata, tt.wantMetadata) {
				t.Errorf("generateGuessWordMetadata() = %v, want %v", gotMetadata, tt.wantMetadata)
			}
		})
	}
}
