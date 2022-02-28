package wengine

import (
	"reflect"
	"testing"
)

func Test_generateWordMetadata(t *testing.T) {
	type args struct {
		guess  string
		secret string
	}
	wantChar := make(map[string]CharMetadata)
	wantChar["w"] = CharMetadata{
		Char:                "w",
		CountInGuess:        1,
		CountInSecret:       1,
		IndexesInGuess:      []int{0},
		IndexesInSecret:     []int{1},
		IndexesValidGuess:   nil,
		IndexesInvalidGuess: []int{0},
	}
	wantChar["x"] = CharMetadata{
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
		wantMetadata WordMetadata
	}{
		{
			name: "Genrate metadata",
			args: args{
				guess:  "wxxxx",
				secret: "swill",
			},
			wantMetadata: WordMetadata{
				Chars: wantChar,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMetadata WordMetadata
			gotMetadata.GenerateWordMetadata(tt.args.guess, tt.args.secret)
			if !reflect.DeepEqual(gotMetadata, tt.wantMetadata) {
				t.Errorf("generateWordMetadata() = %v, want %v", gotMetadata, tt.wantMetadata)
			}
		})
	}
}

func TestCharMetadata_FoundAllSecretChar(t *testing.T) {
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
			cm := &CharMetadata{
				Char:                tt.fields.Char,
				CountInGuess:        tt.fields.CountInGuess,
				CountInSecret:       tt.fields.CountInSecret,
				IndexesInGuess:      tt.fields.IndexesInGuess,
				IndexesInSecret:     tt.fields.IndexesInSecret,
				IndexesValidGuess:   tt.fields.IndexesValidGuess,
				IndexesInvalidGuess: tt.fields.IndexesInvalidGuess,
			}
			if got := cm.FoundAllSecretChar(); got != tt.want {
				t.Errorf("CharMetadata.FoundAllSecretChar() = %v, want %v", got, tt.want)
			}
		})
	}
}
