package wengine

import (
	"sort"
	"testing"
)

type fields struct {
	Words map[int][]string
}
type args struct {
	length int
}

var lenThree = []string{
	"the",
	"its",
	"one",
}
var lenFour = []string{
	"this",
	"word",
	"four",
}
var wordList = map[int][]string {
	3: lenThree,
	4: lenFour,
}

func TestWordList_GetRandomWord(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "Get random word with len 3",
			fields: fields{
				Words: wordList,
			},
			args: args{
				length: 3,
			},
			want: lenThree,
		},
		{
			name: "Get random word with len 4",
			fields: fields{
				Words: wordList,
			},
			args: args{
				length: 4,
			},
			want: lenFour,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wl := &WordList{
				Words: tt.fields.Words,
			}
			got := wl.GetRandomWord(tt.args.length)
			sort.Strings(tt.want)
			searchIdx := sort.SearchStrings(tt.want, got)
			if searchIdx == len(tt.want) || tt.want[searchIdx] != got {
				t.Errorf("WordList.GetRandomWord() = %v, want one of %v", got, tt.want)
			}
		})
	}
}

func TestWordList_HasWord(t *testing.T) {
	type fields struct {
		Words          map[int][]string
	}
	type args struct {
		word string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "HasWord will return true",
			fields: fields{
				Words: wordList,
			},
			args: args{
				word: "one",
			},
			want: true,
		},
		{
			name: "HasWord will return false",
			fields: fields{
				Words: wordList,
			},
			args: args{
				word: "five",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wl := &WordList{
				Words:          tt.fields.Words,
			}
			if got := wl.HasWord(tt.args.word); got != tt.want {
				t.Errorf("WordList.HasWord() = %v, want %v", got, tt.want)
			}
		})
	}
}
