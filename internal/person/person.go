// Package person can perform various operations on person data.
package person

import (
	"fmt"
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
	const needNumFields = 5
	if l := len(fields); l != needNumFields {
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
