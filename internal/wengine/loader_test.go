package wengine

import (
	"embed"
	"reflect"
	"testing"

	"golang.org/x/text/language"
)

//go:embed test-mocks
var testFsLoader embed.FS

func Test_loadWordList(t *testing.T) {
	type args struct {
		wfs WordFileSystem
	}
	tests := []struct {
		name         string
		args         args
		wantWordList WordList
	}{
		{
			name: "Test word loader",
			args: args{
				wfs: WordFileSystem{
					fs: testFsLoader,
					validFilePathTemplate: "test-mocks/%s/valid",
					invalidFilePathTemplate: "test-mocks/%s/invalid",
					locale: language.English,
				},
			},
			wantWordList: WordList{
				Words: map[int][]string{
					2: {
						"hi",
					},
					4: {
						"test",
					},
					5: {
						"hello",
						"valid",
					},
					8: {
						"language",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWordList := loadWordList(tt.args.wfs); !reflect.DeepEqual(gotWordList, tt.wantWordList) {
				t.Errorf("loadWordList() = %v, want %v", gotWordList, tt.wantWordList)
			}
		})
	}
}
