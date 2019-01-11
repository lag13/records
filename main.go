package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/lag13/records/internal/multicsv"
	"github.com/lag13/records/internal/person"
)

func prependFileInfo(fileName string, lineNum int, msgs []string) []string {
	info := []string{}
	for _, msg := range msgs {
		if lineNum == 0 {
			info = append(info, fmt.Sprintf("%s:%s", fileName, msg))
			continue
		}
		info = append(info, fmt.Sprintf("%s:%d: %s", fileName, lineNum, msg))
	}
	return info
}

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
	filesRecords := [][][]string{}
	{ // validate the syntax of the data
		for _, fh := range fhs {
			records, csvParseErrs := multicsv.ReadAll("|, ", 5, fh)
			if len(csvParseErrs) > 0 {
				parseErrs = append(parseErrs, prependFileInfo(fh.Name(), 0, csvParseErrs)...)
				continue
			}
			filesRecords = append(filesRecords, records)
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
			for j, line := range filesRecords[i] {
				p, semParseErrs := person.Parse(line)
				if len(semParseErrs) > 0 {
					parseErrs = append(parseErrs, prependFileInfo(fileName, j+1, semParseErrs)...)
					continue
				}
				persons = append(persons, p)
			}
		}
	}
	return persons, parseErrs
}

const defaultSort = "gender-lastname-asc"

var sortStyleToSortFn = map[string]func([]person.Person){
	defaultSort:     person.SortGenderLastNameAsc,
	"birthdate-asc": person.SortBirthdateAsc,
	"lastname-desc": person.SortLastNameDesc,
}

type sortStyle struct {
	str string
	fn  func([]person.Person)
}

func (s sortStyle) String() string {
	return s.str
}

func (s *sortStyle) Set(str string) error {
	sortFn, ok := sortStyleToSortFn[str]
	if !ok {
		possibleSortStyles := []string{}
		for key := range sortStyleToSortFn {
			possibleSortStyles = append(possibleSortStyles, key)
		}
		sort.Strings(possibleSortStyles)
		return fmt.Errorf("invalid value, allowed values are %s", strings.Join(possibleSortStyles, ", "))
	}
	s.str = str
	s.fn = sortFn
	return nil
}

func main() {
	var ss = sortStyle{str: defaultSort, fn: sortStyleToSortFn[defaultSort]}
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.Var(&ss, "sort", "how to sort the data")
	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}
	persons, errs := parseDataFromFiles(fs.Args())
	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr, strings.Join(errs, "\n"))
		os.Exit(1)
	}
	ss.fn(persons)
	for _, p := range persons {
		fmt.Println(person.Marshal(p))
	}
}
