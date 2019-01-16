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

func TestParse(t *testing.T) {
	tests := []struct {
		name               string
		s                  string
		delimiters         string
		numFieldsPerRecord int
		wantRecord         []string
		wantParseErr       string
	}{
		{
			name:               "no delimiters in string",
			s:                  "heytherearenoseparators",
			delimiters:         "|, ",
			numFieldsPerRecord: 0,
			wantRecord:         nil,
			wantParseErr:       "there are no delimiters",
		},
		{
			name:               "multiple delimiters in string",
			s:                  "hey|there,buddy !",
			delimiters:         "|, ",
			numFieldsPerRecord: 3,
			wantRecord:         nil,
			wantParseErr:       "there should only be one type of separator but multiple ('|', ',', ' ') were specified",
		},
		{
			name:               "incorrect number of fields in record",
			s:                  "hey,there,buddy",
			delimiters:         "|, ",
			numFieldsPerRecord: 7,
			wantRecord:         nil,
			wantParseErr:       "there were 3 fields when there should have been 7",
		},
		{
			name:               "incorrect number of fields in record",
			s:                  "hey&there&buddy",
			delimiters:         "|, &",
			numFieldsPerRecord: 3,
			wantRecord:         []string{"hey", "there", "buddy"},
			wantParseErr:       "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			record, parseErr := multicsv.Parse(test.s, test.delimiters, test.numFieldsPerRecord)
			if got, want := record, test.wantRecord; !reflect.DeepEqual(got, want) {
				t.Errorf("got record %+v, want %+v", got, want)
			}
			if got, want := parseErr, test.wantParseErr; got != want {
				t.Errorf("got parse error %q, want %q", got, want)
			}
		})
	}
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
				"1: there should only be one type of separator but multiple ('|', ',', ' ') were specified",
				"2: there are no delimiters",
				"4: there were 3 fields when there should have been 5",
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
