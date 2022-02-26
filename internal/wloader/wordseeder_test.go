package wloader

import (
	"embed"
	"reflect"
	"testing"

	"github.com/tanmancan/gwordle/v1/internal/dictengine"
)

//go:embed static/*
var testFs embed.FS

func Test_loadWordList(t *testing.T) {
	type args struct {
		language string
		fs       embed.FS
	}
	tests := []struct {
		name         string
		args         args
		wantWordList dictengine.WordList
	}{
		{
			name: "Test word loader",
			args: args{
				language: "test",
				fs: testFs,
			},
			wantWordList: dictengine.WordList{
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
			if gotWordList := loadWordList(tt.args.language, tt.args.fs); !reflect.DeepEqual(gotWordList, tt.wantWordList) {
				t.Errorf("loadWordList() = %v, want %v", gotWordList, tt.wantWordList)
			}
		})
	}
}
