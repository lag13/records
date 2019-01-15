package person_test

import (
	"fmt"
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

func TestMarshal(t *testing.T) {
	tests := []struct {
		p       person.Person
		wantStr string
	}{
		{
			person.Person{"Last", "First", "Gender", "Color", time.Date(2003, time.May, 15, 0, 0, 0, 0, time.UTC)},
			"Last,First,Gender,Color,05/15/2003",
		},
		{
			person.Person{"Bobbo", "Bob", "Male", "Grey", time.Date(1998, time.December, 2, 0, 0, 0, 0, time.UTC)},
			"Bobbo,Bob,Male,Grey,12/02/1998",
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("running test %d", i), func(t *testing.T) {
			if got, want := person.Marshal(test.p), test.wantStr; got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}

func TestSorts(t *testing.T) {
	tests := []struct {
		sortFn      func([]person.Person)
		persons     []person.Person
		wantPersons []person.Person
	}{
		{
			person.SortGenderLastNameAsc,
			[]person.Person{
				{LastName: "Aarons", Gender: "Male"},
				{LastName: "Brady", Gender: "Female"},
				{LastName: "Aarons", Gender: "Female"},
				{LastName: "anderson", Gender: "Female"},
				{LastName: "Zed", Gender: "Female"},
				{LastName: "Tom", Gender: "Male"},
				{LastName: "Bob", Gender: "Male"},
			},
			[]person.Person{
				{LastName: "Aarons", Gender: "Female"},
				{LastName: "anderson", Gender: "Female"},
				{LastName: "Brady", Gender: "Female"},
				{LastName: "Zed", Gender: "Female"},
				{LastName: "Aarons", Gender: "Male"},
				{LastName: "Bob", Gender: "Male"},
				{LastName: "Tom", Gender: "Male"},
			},
		},
		{
			person.SortBirthdateAsc,
			[]person.Person{
				{DateOfBirth: time.Date(1900, time.December, 2, 0, 0, 0, 0, time.UTC)},
				{DateOfBirth: time.Date(2000, time.December, 2, 0, 0, 0, 0, time.UTC)},
				{DateOfBirth: time.Date(1998, time.May, 2, 0, 0, 0, 0, time.UTC)},
				{DateOfBirth: time.Date(1998, time.December, 2, 0, 0, 0, 0, time.UTC)},
			},
			[]person.Person{
				{DateOfBirth: time.Date(1900, time.December, 2, 0, 0, 0, 0, time.UTC)},
				{DateOfBirth: time.Date(1998, time.May, 2, 0, 0, 0, 0, time.UTC)},
				{DateOfBirth: time.Date(1998, time.December, 2, 0, 0, 0, 0, time.UTC)},
				{DateOfBirth: time.Date(2000, time.December, 2, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			person.SortLastNameDesc,
			[]person.Person{
				{LastName: "Aarons"},
				{LastName: "Brady"},
				{LastName: "anderson"},
				{LastName: "Zed"},
				{LastName: "Tom"},
				{LastName: "Bob"},
			},
			[]person.Person{
				{LastName: "Zed"},
				{LastName: "Tom"},
				{LastName: "Brady"},
				{LastName: "Bob"},
				{LastName: "anderson"},
				{LastName: "Aarons"},
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("testing sorting function %d", i), func(t *testing.T) {
			test.sortFn(test.persons)
			if got, want := test.persons, test.wantPersons; !reflect.DeepEqual(got, want) {
				// TODO: This error message if the test fails is simply
				// terrible but it's such an easy test to pass that I don't
				// care right now. I would really like to go back and make the
				// comparisons better though.
				t.Errorf("got sorted list %v, wanted %v", got, want)
			}
		})
	}
}
