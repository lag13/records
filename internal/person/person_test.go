package person_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/lag13/records/internal/person"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		fields     []string
		wantPerson person.Person
		parseErrs  []string
	}{
		{
			name:      "invalid number of fields",
			fields:    []string{"one field for some reason"},
			parseErrs: []string{"fields slice has length 1 but must be exactly length 5"},
		},
		{
			// TODO: The error message for this test (and
			// others like the e2e ones) are really bad
			// because it's tough for the user to see
			// clearly where the mismatch between got and
			// want is. Try to improve this.
			name:   "invalid fields",
			fields: []string{"", "", "", "", "2019"},
			parseErrs: []string{
				"last name (field 1) must be a non-empty string",
				"first name (field 2) must be a non-empty string",
				"gender (field 3) must be a non-empty string",
				"favorite color (field 4) must be a non-empty string",
				"date of birth (field 5) must have the format YYYY-MM-DD",
			},
		},
		{
			name:       "valid fields",
			fields:     []string{"Last", "First", "Gender", "Color", "2006-04-17"},
			wantPerson: person.Person{"Last", "First", "Gender", "Color", time.Date(2006, 4, 17, 0, 0, 0, 0, time.UTC)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotPerson, parseErrs := person.Parse(test.fields)
			if got, want := parseErrs, test.parseErrs; !reflect.DeepEqual(got, want) {
				t.Errorf("got error %v, want %v", got, want)
			}
			if got, want := gotPerson, test.wantPerson; !reflect.DeepEqual(got, want) {
				t.Errorf("got person %+v, want %+v", got, want)
			}
		})
	}
}
