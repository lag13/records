package multicsv_test

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/lag13/records/internal/multicsv"
)

type mockErrReader struct {
}

func (m mockErrReader) Read([]byte) (int, error) {
	return 0, errors.New("non-nil error")
}

func TestReadAll(t *testing.T) {
	tests := []struct {
		name               string
		delimiters         string
		numFieldsPerRecord int
		content            io.Reader
		wantFields         [][]string
		wantParseErrs      []string
	}{
		{
			name:               "so many problems with the file",
			delimiters:         "|, ",
			numFieldsPerRecord: 5,
			content: strings.NewReader(`hey|there,buddy !
noseps
this,is,an,okay,line
too|few|seps`),
			wantParseErrs: []string{
				"1: there should only be one type of separator in a single line but multiple separators ('|', ',', ' ') were specified",
				"2: there is only one field in the record but there should be 5",
				"4: there were only 3 fields when there should have been 5",
			},
		},
		{
			name:               "return the fields",
			delimiters:         "|, &",
			numFieldsPerRecord: 5,
			content: strings.NewReader(`one|two|three|four|five
6,7,8,,10
11 12 13 14 15
16&17&18&19&20`),
			wantFields: [][]string{
				{"one", "two", "three", "four", "five"},
				{"6", "7", "8", "", "10"},
				{"11", "12", "13", "14", "15"},
				{"16", "17", "18", "19", "20"},
			},
			wantParseErrs: nil,
		},
		{
			name:          "error when reading file",
			content:       mockErrReader{},
			wantParseErrs: []string{`unexpected error reading file: non-nil error`},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields, parseErrs := multicsv.ReadAll(test.delimiters, test.numFieldsPerRecord, test.content)
			if got, want := fields, test.wantFields; !reflect.DeepEqual(got, want) {
				t.Errorf("got fields %+v, want %+v", got, want)
			}
			if got, want := parseErrs, test.wantParseErrs; !reflect.DeepEqual(got, want) {
				t.Errorf("got parse errs %+v, want %+v", got, want)
			}
		})
	}
}
