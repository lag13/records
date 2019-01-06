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

// ParseFile converts file information into information which can be
// analyzed as well as returning errors from parsing. TODO: Right now
// it just returns errors from parsing, we still need to add the part
// where it returns useful data.
func ParseFile(file File) ([][]string, string) {
	if file.OpenErr != nil {
		if os.IsNotExist(file.OpenErr) {
			return nil, fmt.Sprintf("%s: file does not exist", file.Name)
		}
		if os.IsPermission(file.OpenErr) {
			return nil, fmt.Sprintf("%s: do not have permission to open this file", file.Name)
		}
		return nil, fmt.Sprintf("%s: encountered an unknown error when opening this file: %v", file.Name, file.OpenErr)
	}
	parseErrs := []string{}
	scanner := bufio.NewScanner(file.Content)
	lineNo := 0
	allFields := [][]string{}
	for scanner.Scan() {
		// TODO: bufio.Scanner will NOT ignore empty lines.
		// I'll keep it that way for now but I wonder if in
		// the future we want to ignore/allow empty lines.
		lineNo++
		seps := whichSeparatorsUsedInLine(scanner.Text())
		if len(seps) == 0 {
			parseErrs = append(parseErrs, fmt.Sprintf("%s:%d: there is only one field in the record but there should be 5", file.Name, lineNo))
			continue
		}
		if len(seps) > 1 {
			sepsStr := fmt.Sprintf("'%c'", seps[0])
			for _, sep := range seps[1:] {
				sepsStr = fmt.Sprintf("%s, '%c'", sepsStr, sep)
			}
			parseErrs = append(parseErrs, fmt.Sprintf("%s:%d: there should only be one type of separator in a single line but multiple separators (%s) were specified", file.Name, lineNo, sepsStr))
			continue
		}
		fields := strings.Split(scanner.Text(), string(seps[0]))
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
		// reminder of a potential refactor.
		const desiredNumFields = 5
		if numFields := len(fields); numFields != desiredNumFields {
			parseErrs = append(parseErrs, fmt.Sprintf("%s:%d: there were only %d fields when there should have been %d", file.Name, lineNo, numFields, desiredNumFields))
			continue
		}
		allFields = append(allFields, fields)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Sprintf("AHHHHH WE SHOULD PROBABLY TEST THIS!!!")
	}
	if len(parseErrs) > 0 {
		allFields = nil
	}
	return allFields, strings.Join(parseErrs, "\n")
}
