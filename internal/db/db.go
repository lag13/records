// Package db contains an in memory database for the API portion of
// this problem. In other words, it just defines a slice which other
// things in this repository can reference.
package db

import "github.com/lag13/records/internal/person"

// Persons is the data we store in the API
var Persons = []person.Person{}
