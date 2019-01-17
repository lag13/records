// Package response defines a response type that will get returned
// from most handlers in this repository.
package response

import "github.com/lag13/records/internal/person"

// Structured is a http response with a structured body.
type Structured struct {
	StatusCode int             `json:"-"`
	Data       []person.Person `json:"data,omitempty"`
	Errors     []string        `json:"errors,omitempty"`
}
