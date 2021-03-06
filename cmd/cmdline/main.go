package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/lag13/records/internal/multicsv"
	"github.com/lag13/records/internal/person"
)

// TODO: I feel like calls to this function should happen in one
// single place (which would also remove the need for said function).
// Like every file which returns data about parse errors should,
// instead of returning type []string, should be returning some sort
// of struct like []struct{fileName, errLines: []struct{lineNum,
// errMsg}} then main will assemble that information in one single
// place. Not sure yet.
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

// TODO: This logic feels too complicated especially for main. Maybe I
// just need to put more of the loops into the units. God I wish Go
// had map and other such operations which operate on collections.
func parseDataFromFiles(fileNames []string) ([]person.Person, []string) {
	type simpleFile struct {
		Name    string
		Content io.Reader
	}
	files := []simpleFile{}
	parseErrs := []string{}
	{ // open all specified files
		for _, fileName := range fileNames {
			fh, err := os.Open(fileName)
			if err != nil {
				parseErrs = append(parseErrs, fmt.Sprint(err))
				continue
			}
			// we cannot do a simple 'defer fh.Close()'
			// because the errcheck tool will say we're
			// forgetting to check an error. Ignoring the
			// error is fine in this case because we are
			// reading from the file as opposed to
			// writing. TODO: check out
			// https://github.com/alecthomas/gometalinter
			// as a way to consolidate static checks AND
			// it seems you can add a comment instructing
			// gometalinter not to run, say errcheck, on a
			// specific line
			defer func(fh *os.File) { _ = fh.Close() }(fh)
			files = append(files, simpleFile{Name: fileName, Content: fh})
		}
	}
	if len(parseErrs) > 0 {
		return nil, parseErrs
	}
	if len(files) == 0 {
		files = append(files, simpleFile{Name: "(standard input)", Content: os.Stdin})
	}
	filesRecords := [][][]string{}
	{ // validate the syntax of the data
		const possibleDelimiters = "|, "
		const numFieldsInRecord = 5
		for _, file := range files {
			records, csvParseErrs := multicsv.ReadAll(possibleDelimiters, numFieldsInRecord, file.Content)
			if len(csvParseErrs) > 0 {
				parseErrs = append(parseErrs, prependFileInfo(file.Name, 0, csvParseErrs)...)
				continue
			}
			filesRecords = append(filesRecords, records)
		}
	}
	// TODO: I don't like the repetition of this check but I'm not
	// sure what to do about it. Hypothetically I could stop every
	// time I encounter any error instead of gathering them but I
	// think its useful to the user if all problems of a certain
	// type (like grammar vs semantic) are printed out at once.
	// Hypothetically I could also not even bother checking
	// parseErrs in this function at all and have it keep
	// accumulating but that would mean some unecessary
	// computation being done for any files without any problems.
	if len(parseErrs) > 0 {
		return nil, parseErrs
	}
	persons := []person.Person{}
	{ // parse each file into structured data which can be sorted
		for i, file := range files {
			for j, line := range filesRecords[i] {
				p, semParseErrs := person.Parse(line)
				if len(semParseErrs) > 0 {
					parseErrs = append(parseErrs, prependFileInfo(file.Name, j+1, semParseErrs)...)
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
