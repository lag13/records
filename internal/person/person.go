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

// These constants are used to dispatch to a certain way of sorting
// the person data.
const (
	GenderThenLastNameAsc = "gender-lastname-asc"
	BirthDateAsc          = "birthdate-asc"
	LastNameDesc          = "lastname-desc"
)

// SortStyles contains all the possible sort styles that can be
// achieved. TODO: As it stands now this variable only gets used in
// main to display all possible values for a particular flag AND it
// has no direct relation to the Sort() function below. This is no
// good because we could (hypothetically) add another item to this
// list, not update Sort(), user thinks they can pass another sort
// type, and then Sort() will fail because it does not actually
// support said type. This should be remedied. I'm picturing this
// variable being a map from sort type to function which sorts a list
// of people. Then where main does the flag validation it can
// literally check if the key exists in this map and if it does not
// then they print the error message otherwise they use the function
// at that key to perform the sorting. Doing it like this also means
// that we don't need the panic inside the Sort() function, which I'm
// not a huge fan of (because there won't be one big function anymore,
// there will be 3 functions which do a particular kind of sorting!).
// I'm glad I thought this through, I knew the panic made me feel a
// bit weird but I wasn't sure as to why.
var SortStyles = []string{GenderThenLastNameAsc, BirthDateAsc, LastNameDesc}

// Sort sorts a slice of Person structs in place in a couple different
// ways.
func Sort(sortStyle string, persons []Person) {
	if sortStyle == BirthDateAsc {
		sort.SliceStable(persons, func(i int, j int) bool {
			return persons[i].DateOfBirth.Before(persons[j].DateOfBirth)
		})
		return
	}
	if sortStyle == LastNameDesc {
		sort.SliceStable(persons, func(i int, j int) bool {
			return strings.ToLower(persons[i].LastName) > strings.ToLower(persons[j].LastName)
		})
		return
	}
	if sortStyle == GenderThenLastNameAsc {
		sort.SliceStable(persons, func(i int, j int) bool {
			if persons[i].Gender == persons[j].Gender {
				return strings.ToLower(persons[i].LastName) < strings.ToLower(persons[j].LastName)
			}
			return persons[i].Gender < persons[j].Gender
		})
		return
	}
	panic(fmt.Sprintf("invalid sort style %q passed", sortStyle))
}
