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

// Parse converts a list of fields into a Person struct. It MUST be
// passed a slice of at least 5 otherwise it will panic.
func Parse(fields []string) (Person, []string) {
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

// SortGenderLastNameAsc sorts a slice of Person structs females first
// then by last name ascending.
func SortGenderLastNameAsc(persons []Person) {
	sort.SliceStable(persons, func(i int, j int) bool {
		if persons[i].Gender == persons[j].Gender {
			return strings.ToLower(persons[i].LastName) < strings.ToLower(persons[j].LastName)
		}
		return persons[i].Gender < persons[j].Gender
	})
}

// SortBirthdateAsc sorts a slice of Person structs by birth date.
func SortBirthdateAsc(persons []Person) {
	sort.SliceStable(persons, func(i int, j int) bool {
		return persons[i].DateOfBirth.Before(persons[j].DateOfBirth)
	})
}

// SortLastNameDesc sorts a slice of Person structs by last name
// descending.
func SortLastNameDesc(persons []Person) {
	sort.SliceStable(persons, func(i int, j int) bool {
		return strings.ToLower(persons[i].LastName) > strings.ToLower(persons[j].LastName)
	})
}
