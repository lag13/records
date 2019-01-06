// Package parsefile parses file information so it can be analyzed.
package parsefile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// File contains raw data coming from a call to os.Open().
type File struct {
	Name    string
	Content io.Reader
	OpenErr error
}

func whichSeparatorsUsedInLine(line string) []rune {
	const supportedSeparators = "|, "
	seenSeps := []rune{}
	for _, sep := range supportedSeparators {
		if strings.ContainsRune(line, sep) {
			seenSeps = append(seenSeps, sep)
		}
	}
	return seenSeps
}

func parseErrPrefix(filename string, lineNum int, msg string) string {
	if lineNum == 0 {
		return fmt.Sprintf("%s: %s", filename, msg)
	}
	return fmt.Sprintf("%s:%d: %s", filename, lineNum, msg)
}

// ParseFile converts file information into information which can be
// analyzed as well as returning errors from parsing.
func ParseFile(file File) ([][]string, string) {
	if file.OpenErr != nil {
		if os.IsNotExist(file.OpenErr) {
			return nil, parseErrPrefix(file.Name, 0, "file does not exist")
		}
		if os.IsPermission(file.OpenErr) {
			return nil, parseErrPrefix(file.Name, 0, "do not have permission to open this file")
		}
		return nil, parseErrPrefix(file.Name, 0, fmt.Sprintf("encountered an unknown error when opening this file: %v", file.OpenErr))
	}
	parseErrs := []string{}
	scanner := bufio.NewScanner(file.Content)
	lineNum := 0
	allFields := [][]string{}
	for scanner.Scan() {
		lineNum++
		seps := whichSeparatorsUsedInLine(scanner.Text())
		// TODO: So far with this code all I've been doing is
		// making sure that the file has the expected "shape"
		// without caring about it's contents similar to what
		// this blog goes through:
		// https://dev.to/matthewsj/you-could-have-designed-the-jsondecode-library-2d8.
		// That all makes me think that I could make this a
		// more general package, a sort of "CSV parser which
		// allows different separators in the same file". But
		// it's probably not worth doing it just yet or maybe
		// not at all! Wanted to mention this though as a
		// reminder of a potential refactor where we pass in
		// the desired number of fields and separators.
		const desiredNumFields = 5
		if len(seps) == 0 {
			parseErrs = append(parseErrs, parseErrPrefix(file.Name, lineNum, fmt.Sprintf("there is only one field in the record but there should be %d", desiredNumFields)))
			continue
		}
		if len(seps) > 1 {
			sepsStr := fmt.Sprintf("'%c'", seps[0])
			for _, sep := range seps[1:] {
				sepsStr = fmt.Sprintf("%s, '%c'", sepsStr, sep)
			}
			parseErrs = append(parseErrs, parseErrPrefix(file.Name, lineNum, fmt.Sprintf("there should only be one type of separator in a single line but multiple separators (%s) were specified", sepsStr)))
			continue
		}
		fields := strings.Split(scanner.Text(), string(seps[0]))
		if numFields := len(fields); numFields != desiredNumFields {
			parseErrs = append(parseErrs, fmt.Sprintf("%s:%d: there were only %d fields when there should have been %d", file.Name, lineNum, numFields, desiredNumFields))
			continue
		}
		allFields = append(allFields, fields)
	}
	if err := scanner.Err(); err != nil {
		return nil, parseErrPrefix(file.Name, 0, fmt.Sprintf("unexpected error reading file: %v", err))
	}
	if len(parseErrs) > 0 {
		allFields = nil
	}
	return allFields, strings.Join(parseErrs, "\n")
}
