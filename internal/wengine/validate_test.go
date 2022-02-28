package wengine

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
		{
			name: "Guess word: nnnnn. Secret word: glint",
			args: args{
				guess:  "nnnnn",
				secret: "glint",
			},
			wantResult: ValidationResult{
				Match: false,
				Chars: []CharacterValidationResult{
					{
						Char:   "n",
						Status: InvalidCharacter,
					},
					{
						Char:   "n",
						Status: InvalidCharacter,
					},
					{
						Char:   "n",
						Status: InvalidCharacter,
					},
					{
						Char:   "n",
						Status: ValidPosition,
					},
					{
						Char:   "n",
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
		Char:                "w",
		CountInGuess:        1,
		CountInSecret:       1,
		IndexesInGuess:      []int{0},
		IndexesInSecret:     []int{1},
		IndexesValidGuess:   nil,
		IndexesInvalidGuess: []int{0},
	}
	wantChar["x"] = GuessWordCharMetadata{
		Char:                "x",
		CountInGuess:        4,
		CountInSecret:       0,
		IndexesInGuess:      []int{1, 2, 3, 4},
		IndexesInSecret:     nil,
		IndexesValidGuess:   nil,
		IndexesInvalidGuess: nil,
	}
	tests := []struct {
		name         string
		args         args
		wantMetadata GuessWordMetadata
	}{
		{
			name: "Genrate metadata",
			args: args{
				guess:  "wxxxx",
				secret: "swill",
			},
			wantMetadata: GuessWordMetadata{
				Chars: wantChar,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMetadata GuessWordMetadata
			gotMetadata.GenerateGuessWordMetadata(tt.args.guess, tt.args.secret)
			if !reflect.DeepEqual(gotMetadata, tt.wantMetadata) {
				t.Errorf("generateGuessWordMetadata() = %v, want %v", gotMetadata, tt.wantMetadata)
			}
		})
	}
}

func TestGuessWordCharMetadata_FoundAllSecretChar(t *testing.T) {
	type fields struct {
		Char                string
		CountInGuess        int
		CountInSecret       int
		IndexesInGuess      []int
		IndexesInSecret     []int
		IndexesValidGuess   []int
		IndexesInvalidGuess []int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Returns true",
			fields: fields{
				Char: "x",
				CountInGuess: 2,
				CountInSecret: 2,
				IndexesInGuess: []int{0, 1},
				IndexesInSecret: []int{0, 1},
				IndexesValidGuess: []int{0, 1},
				IndexesInvalidGuess: nil,
			},
			want: true,
		},
		{
			name: "Returns true",
			fields: fields{
				Char: "x",
				CountInGuess: 4,
				CountInSecret: 2,
				IndexesInGuess: []int{0, 1, 4, 7},
				IndexesInSecret: []int{4, 7},
				IndexesValidGuess: []int{4, 7},
				IndexesInvalidGuess: nil,
			},
			want: true,
		},
		{
			name: "Returns false",
			fields: fields{
				Char: "x",
				CountInGuess: 2,
				CountInSecret: 2,
				IndexesInGuess: []int{0, 1},
				IndexesInSecret: []int{2, 3},
				IndexesValidGuess: []int{0, 1},
				IndexesInvalidGuess: nil,
			},
			want: false,
		},
		{
			name: "Returns false",
			fields: fields{
				Char: "x",
				CountInGuess: 4,
				CountInSecret: 3,
				IndexesInGuess: []int{0, 1, 4, 7},
				IndexesInSecret: []int{4, 7, 8},
				IndexesValidGuess: []int{4, 7},
				IndexesInvalidGuess: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &GuessWordCharMetadata{
				Char:                tt.fields.Char,
				CountInGuess:        tt.fields.CountInGuess,
				CountInSecret:       tt.fields.CountInSecret,
				IndexesInGuess:      tt.fields.IndexesInGuess,
				IndexesInSecret:     tt.fields.IndexesInSecret,
				IndexesValidGuess:   tt.fields.IndexesValidGuess,
				IndexesInvalidGuess: tt.fields.IndexesInvalidGuess,
			}
			if got := cm.FoundAllSecretChar(); got != tt.want {
				t.Errorf("GuessWordCharMetadata.FoundAllSecretChar() = %v, want %v", got, tt.want)
			}
		})
	}
}
