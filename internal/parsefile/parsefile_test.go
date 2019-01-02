package parsefile_test

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/lag13/records/internal/parsefile"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name         string
		file         parsefile.File
		wantParseErr string
	}{
		{
			name: "file does not exist",
			file: parsefile.File{
				Name:    "some/file.txt",
				Content: nil,
				OpenErr: os.ErrNotExist,
			},
			wantParseErr: "some/file.txt: file does not exist",
		},
		{
			name: "file has bad permissions",
			file: parsefile.File{
				Name:    "badperm/file.txt",
				Content: nil,
				OpenErr: os.ErrPermission,
			},
			wantParseErr: "badperm/file.txt: do not have permission to open this file",
		},
		{
			name: "unknown error",
			file: parsefile.File{
				Name:    "itsa/mystery.txt",
				Content: nil,
				OpenErr: errors.New("unknown error"),
			},
			wantParseErr: "itsa/mystery.txt: encountered an unknown error when opening this file: unknown error",
		},
		{
			name: "unknown open error",
			file: parsefile.File{
				Name:    "itsa/mystery.txt",
				Content: nil,
				OpenErr: errors.New("unknown error"),
			},
			wantParseErr: "itsa/mystery.txt: encountered an unknown error when opening this file: unknown error",
		},
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := parsefile.ParseFile(test.file), test.wantParseErr; !reflect.DeepEqual(got, want) {
				t.Errorf("got parse err %q, want %q", got, want)
			}
		})
	}
}
