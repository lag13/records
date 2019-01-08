// Package person can perform various operations on person data.
package person

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Person contains data about a person.
type Person struct {
	LastName      string
	FirstName     string
	Gender        string
	FavoriteColor string
	DateOfBirth   time.Time
}

// Parse converts a list of fields into a Person struct.
func Parse(fields []string) (Person, []string) {
	// TODO: This constant 5 now lives in 2 places which seems not
	// good. I think it should be passed into this package and the
	// parsefile package. That would mean making these 2 packages
	// more general.
	const needNumFields = 5
	if l := len(fields); l != needNumFields {
		// TODO: I think this should be a panic instead of a
		// parse error thingy (we could even just let this
		// function panic naturally). Because there is NEVER a
		// reason this function should be passed a slice that
		// has a length less than 5. If that ever happens it
		// is clearly the callers fault so why should this
		// function worry about checking for a condition which
		// should never happen?
		return Person{}, []string{fmt.Sprintf("fields slice has length %d but must be exactly length %d", l, needNumFields)}
	}
	parseErrs := []string{}
	nonEmptyFieldNames := []string{"last name", "first name", "gender", "favorite color"}
	for i, fieldName := range nonEmptyFieldNames {
		if fields[i] != "" {
			continue
		}
		parseErrs = append(parseErrs, fmt.Sprintf("%s (field %d) must be a non-empty string", fieldName, i+1))
	}
	// https://stackoverflow.com/questions/14106541/go-parsing-date-time-strings-which-are-not-standard-formats
	layout := "2006-01-02"
	dob, err := time.Parse(layout, fields[4])
	if err != nil {
		parseErrs = append(parseErrs, "date of birth (field 5) must have the format YYYY-MM-DD")
	}
	if len(parseErrs) > 0 {
		return Person{}, parseErrs
	}
	return Person{fields[0], fields[1], fields[2], fields[3], dob}, nil
}

// Marshal converts a Person struct into a CSV row
func Marshal(p Person) string {
	year, month, day := p.DateOfBirth.Date()
	return fmt.Sprintf("%s,%s,%s,%s,%02d/%02d/%d", p.LastName, p.FirstName, p.Gender, p.FavoriteColor, month, day, year)
}

// SortByGenderThenLastName sorts a slice of Person structs in place
// oby gender then by last name.
func SortByGenderThenLastName(persons []Person) {
	sort.SliceStable(persons, func(i int, j int) bool {
		if persons[i].Gender == persons[j].Gender {
			return strings.ToLower(persons[i].LastName) < strings.ToLower(persons[j].LastName)
		}
		return persons[i].Gender < persons[j].Gender
	})
}
