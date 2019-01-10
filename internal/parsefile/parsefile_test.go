package parsefile_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/lag13/records/internal/parsefile"
)

type mockErrReader struct {
}

func (m mockErrReader) Read([]byte) (int, error) {
	return 0, errors.New("non-nil error")
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		name         string
		file         parsefile.File
		wantFields   [][]string
		wantParseErr string
	}{
		{
			name: "so many problems with the file",
			file: parsefile.File{
				Name: "hey.txt",
				Content: strings.NewReader(`hey|there,buddy !
noseps
this,is,an,okay,line
too|few|seps`),
			},
			wantParseErr: `hey.txt:1: there should only be one type of separator in a single line but multiple separators ('|', ',', ' ') were specified
hey.txt:2: there is only one field in the record but there should be 5
hey.txt:4: there were only 3 fields when there should have been 5`,
		},
		{
			name: "return the fields",
			file: parsefile.File{
				Name: "hey.txt",
				Content: strings.NewReader(`one|two|three|four|five
6,7,8,,10
11 12 13 14 15`),
			},
			wantFields: [][]string{
				{"one", "two", "three", "four", "five"},
				{"6", "7", "8", "", "10"},
				{"11", "12", "13", "14", "15"},
			},
			wantParseErr: "",
		},
		{
			name: "error when reading file",
			file: parsefile.File{
				Name:    "err.txt",
				Content: mockErrReader{},
			},
			wantParseErr: `err.txt: unexpected error reading file: non-nil error`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields, parseErr := parsefile.ParseFile(test.file)
			if got, want := fields, test.wantFields; !reflect.DeepEqual(got, want) {
				t.Errorf("got fields %+v, want %+v", got, want)
			}
			if got, want := parseErr, test.wantParseErr; got != want {
				t.Errorf("got parse err %q, want %q", got, want)
			}
		})
	}
}
