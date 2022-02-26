package fhandler

import (
	"bufio"
	"embed"
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

//go:embed static/test_words.txt
var testFs embed.FS
func TestWordListFileReader(t *testing.T) {
	type args struct {
		path string
		fs embed.FS
	}
	tests := []struct {
		name string
		args args
		want *strings.Reader
		wantErr bool
	}{
		{
			name: "Test file reader",
			args: args{
				path: "static/test_words.txt",
				fs: testFs,
			},
			want: strings.NewReader("testing\nvarious\nwords\nlist\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WordListFileReader(tt.args.path, tt.args.fs);
			if err != nil {
				t.Errorf("WordListFileReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WordListFileReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
