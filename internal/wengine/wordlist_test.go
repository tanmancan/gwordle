package wengine

import (
	"sort"
	"testing"
)

type fields struct {
	Words map[int][]string
	FilterWords []string
}
type args struct {
	length int
}

var lenThree = []string{
	"the",
	"its",
	"one",
	"two",
	"and",
	"for",
	"pun",
	"cat",
	"dog",
	"fat",
	"but",
	"not",
}
var lenFour = []string{
	"this",
	"word",
	"four",
	"nine",
	"nice",
	"port",
	"face",
	"just",
	"zero",
	"nope",
}
var lenFive = []string{
	"juice",
	"world",
}
var filterFive = []string{
	"juice",
	"timer",
	"funny",
}
var lenSix = []string{
	"filter",
	"potato",
}
var filterSix = []string{
	"potato",
	"relief",
	"crayon",
}
var wordList = map[int][]string{
	3: lenThree,
	4: lenFour,
	5: lenFive,
	6: lenSix,
}

func TestWordList_GetRandomWord(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "Return random word with len 3",
			fields: fields{
				Words: wordList,
			},
			args: args{
				length: 3,
			},
			want: lenThree,
		},
		{
			name: "Return random word with len 4",
			fields: fields{
				Words: wordList,
			},
			args: args{
				length: 4,
			},
			want: lenFour,
		},
		{
			name: "Return filtered word with len 5",
			fields: fields{
				Words: wordList,
				FilterWords: filterFive,
			},
			args: args{
				length: 5,
			},
			want: []string{
				"world",
			},
		},
		{
			name: "Return filtered word with len 6",
			fields: fields{
				Words: wordList,
				FilterWords: filterSix,
			},
			args: args{
				length: 6,
			},
			want: []string{
				"filter",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wl := &WordList{
				Words: tt.fields.Words,
				FilterWords: tt.fields.FilterWords,
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
		Words map[int][]string
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
				Words: tt.fields.Words,
			}
			if got := wl.HasWord(tt.args.word); got != tt.want {
				t.Errorf("WordList.HasWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWordList_SetFilterWord(t *testing.T) {
	type fields struct {
		FilterWords []string
	}
	type args struct {
		word string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Adds a word to filter list",
			fields: fields{
				FilterWords: []string{},
			},
			args: args{
				word: "testingfilter",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wl := &WordList{
				FilterWords: tt.fields.FilterWords,
			}
			if got := wl.HasFilterWord(tt.args.word); got {
				t.Errorf("WordList.SetFilterWord() = %v, want %v", got, false)
			}
			wl.SetFilterWord(tt.args.word)
			if got := wl.HasFilterWord(tt.args.word); !got {
				t.Errorf("WordList.SetFilterWord() = %v, want %v", got, true)
			}
		})
	}
}
