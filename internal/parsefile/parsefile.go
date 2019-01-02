// Package parsefile parses file information so it can be analyzed.
package parsefile

import "io"

// File contains raw data coming from a call to os.Open().
type File struct {
	Name    string
	Content io.Reader
	OpenErr error
}

// ParseFile converts file information into information which can be
// analyzed as well as returning errors from parsing. TODO: Right now
// it just returns errors from parsing, we still need to add the part
// where it returns useful data.
func ParseFile(file File) string {
	return ""
}
