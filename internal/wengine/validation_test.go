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
				Chars: []CharValidationResult{
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
				Chars: []CharValidationResult{
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
				Chars: []CharValidationResult{
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
				Chars: []CharValidationResult{
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
				Chars: []CharValidationResult{
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
				Chars: []CharValidationResult{
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
		{
			name: "Guess word: nnnnn. Secret word: reccuring",
			args: args{
				guess:  "ccccccick",
				secret: "reccuring",
			},
			wantResult: ValidationResult{
				Match: false,
				Chars: []CharValidationResult{
					{
						Char:   "c",
						Status: InvalidCharacter,
					},
					{
						Char:   "c",
						Status: InvalidCharacter,
					},
					{
						Char:   "c",
						Status: ValidPosition,
					},
					{
						Char:   "c",
						Status: ValidPosition,
					},
					{
						Char:   "c",
						Status: InvalidCharacter,
					},
					{
						Char:   "c",
						Status: InvalidCharacter,
					},
					{
						Char:   "i",
						Status: ValidPosition,
					},
					{
						Char:   "c",
						Status: InvalidCharacter,
					},
					{
						Char:   "k",
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
