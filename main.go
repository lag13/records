package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lag13/records/internal/parsefile"
	"github.com/lag13/records/internal/person"
)

// TODO: Although I like this logic better because it is a little more
// pipeline'y and because it prints out similar groupings of errors
// (file related errors, file content structure related errors, and
// file content semantic related errors) it still feels too
// complicated especially for main. I should be able to simplify it
// somehow. Maybe I just need to put more of the loops into the units.
// I think I did not do something like that originally because I felt
// like exposing a "do transformation on a single thing" should be all
// that is necessary for the unit to expose. Applying that
// transformation on multiple things is pretty trivial. Hmmmmmm
func parseDataFromFiles(fileNames []string) ([]person.Person, []string) {
	fhs := []*os.File{}
	parseErrs := []string{}
	{ // open all specified files
		for _, fileName := range fileNames {
			fh, err := os.Open(fileName)
			if err != nil {
				parseErrs = append(parseErrs, fmt.Sprint(err))
				continue
			}
			// cannot do a simple 'defer fh.Close()'
			// because the errcheck tool will say we're
			// forgetting to check an error. Ignoring the
			// error is fine because we are reading from
			// the file. TODO: check out
			// https://github.com/alecthomas/gometalinter
			// as a way to consolidate static checks AND
			// it seems you can add a comment instructing
			// gometalinter not to run, say errcheck, on a
			// specific line
			defer func(fh *os.File) { _ = fh.Close() }(fh)
			fhs = append(fhs, fh)
		}
	}
	if len(parseErrs) > 0 {
		return nil, parseErrs
	}
	rawFileContents := [][][]string{}
	{ // validate the syntax of the data
		for _, fh := range fhs {
			lines, parseErr := parsefile.ParseFile(parsefile.File{
				Name:    fh.Name(),
				Content: fh,
			})
			if parseErr != "" {
				parseErrs = append(parseErrs, parseErr)
				continue
			}
			rawFileContents = append(rawFileContents, lines)
		}
	}
	// TODO: I don't like the repetition of this check but I'm not
	// sure what to do about it. Hypothetically I could stop
	// everytime I encounter any error instead of gathering them
	// but I think its useful to the user if all problems of a
	// certain type (like grammar vs semantic) are printed out at
	// once. Hypothetically I could also not even bother checking
	// parseErrs in this function at all and have it keep
	// accumulating but that would mean some unecessary
	// computation being done for any files without any problems.
	if len(parseErrs) > 0 {
		return nil, parseErrs
	}
	persons := []person.Person{}
	{ // parse each file into structured data which can be sorted
		for i, fileName := range fileNames {
			for j, line := range rawFileContents[i] {
				p, semParseErrs := person.Parse(line)
				if len(semParseErrs) > 0 {
					// TODO: All this indentation
					// is just... awful.
					for _, semParseErr := range semParseErrs {
						// TODO: This format
						// is repeated here
						// and inside
						// parsefile.go and it
						// should be
						// abstracted.
						parseErrs = append(parseErrs, fmt.Sprintf("%s:%d: %s", fileName, j+1, semParseErr))
					}
					continue
				}
				persons = append(persons, p)
			}
		}
	}
	return persons, parseErrs
}

// TODO: Earlier I was planning the person package be responsible for
// knowing the possible values of this sort style string but I think
// that's wrong. Why should that package be dictating valid flag
// arguments? Put another way, if I want to sort something in a
// specific way, I would prefer to just call the function that does it
// for me instead of needing to know which string argument to pass in
// which *gets* me the sort that I want. But now it feels like we're
// adding enough logic where I might feel better making some sort of
// "flags" package. But what would we test exactly? And what would we
// expose? I suppose we'd check that the default values and name are
// as expected? And then we can test any of the custom flag value
// stuff that we do?
type sortStyle string

func (s sortStyle) String() string {
	return string(s)
}

func (s *sortStyle) Set(str string) error {
	if !contains(person.SortStyles, str) {
		return fmt.Errorf("invalid value, allowed values are %s", strings.Join(person.SortStyles, ", "))
	}
	*s = sortStyle(str)
	return nil
}

func contains(xs []string, y string) bool {
	for _, x := range xs {
		if x == y {
			return true
		}
	}
	return false
}

func main() {
	var sortStyle sortStyle = sortStyle(person.SortStyles[0])
	flag.Var(&sortStyle, "sort", "how to sort the data")
	flag.Parse()
	persons, errs := parseDataFromFiles(flag.Args())
	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr, strings.Join(errs, "\n"))
		os.Exit(1)
	}
	person.Sort(string(sortStyle), persons)
	for _, p := range persons {
		fmt.Println(person.Marshal(p))
	}
}
