package fhandler

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const testPath = "test_words.txt"

func TestWordListFileScanner(t *testing.T) {
	mockWordList := []string{
		"testing",
		"various",
		"words",
		"list",
	}
	mockContent := strings.Join(mockWordList, "\n")
	mockReader := strings.NewReader(mockContent)

	type args struct {
		reader *strings.Reader
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: fmt.Sprintf("Test file scanner with mock reader. Reader content: %s", mockContent),
			args: args{
				reader: mockReader,
			},
			want: mockWordList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotWordList []string
			got := WordListFileScanner(tt.args.reader)
			got.Split(bufio.ScanLines)

			for got.Scan() {
				gotWordList = append(gotWordList, got.Text())
			}

			if !reflect.DeepEqual(gotWordList, tt.want) {
				t.Errorf("WordListFileScanner() = %v, want %v", gotWordList, mockWordList)
			}
		})
	}
}

func TestWordListFileReader(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *strings.Reader
	}{
		{
			name: "Test file reader",
			args: args{
				path: "test_words.txt",
			},
			want: strings.NewReader("testing\nvarious\nwords\nlist\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WordListFileReader(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WordListFileReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
