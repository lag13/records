package parsefile_test

import (
	"os"
	"reflect"
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := parsefile.ParseFile(test.file), test.wantParseErr; !reflect.DeepEqual(got, want) {
				t.Errorf("got parse err %q, want %q", got, want)
			}
		})
	}
}
