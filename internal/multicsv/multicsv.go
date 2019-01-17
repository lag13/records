// Package multicsv is similar to go's encoding/csv package except it
// can parse a single file which has multiple different separators but
// besides that it is much less feature rich.
package multicsv

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func whichSeparatorsUsedInLine(line string, delimiters string) []rune {
	seenSeps := []rune{}
	for _, sep := range delimiters {
		if strings.ContainsRune(line, sep) {
			seenSeps = append(seenSeps, sep)
		}
	}
	return seenSeps
}

// Parse converts a string containing a string delimited by something
// and converts it to a []string
func Parse(s string, delimiters string, numFieldsPerRecord int) ([]string, string) {
	seps := whichSeparatorsUsedInLine(s, delimiters)
	// TODO: I feel like this case is unecessary and a little
	// strange since you could hypothetically pass in
	// numFieldsPerRecord = 1 in which case there need not be any
	// delimiters. Maybe I should not have bothered making
	// parameters out of delimiters and numFieldsPerRecord
	if len(seps) == 0 {
		return nil, "there are no delimiters"
	}
	if len(seps) > 1 {
		sepsStr := fmt.Sprintf("'%c'", seps[0])
		for _, sep := range seps[1:] {
			sepsStr = fmt.Sprintf("%s, '%c'", sepsStr, sep)
		}
		return nil, fmt.Sprintf("there should only be one type of separator but multiple (%s) were specified", sepsStr)
	}
	fields := strings.Split(s, string(seps[0]))
	if numFields := len(fields); numFields != numFieldsPerRecord {
		return nil, fmt.Sprintf("there were %d fields when there should have been %d", numFields, numFieldsPerRecord)
	}
	return fields, ""
}

// ReadAll reads all records out of the Reader.
func ReadAll(delimiters string, numFieldsPerRecord int, r io.Reader) ([][]string, []string) {
	parseErrs := []string{}
	// TODO: I could see us wanting to ignore empty lines but
	// bufio.Scanner does NOT ignore empty lines. Keep this in
	// mind. TODO: Keep in mind that bufio.Scanner has a limited
	// buffer size:
	// https://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go.
	// I don't imagine it would be a problem especially for this
	// practice problem but I wanted to think about it a little
	// more if I tried to make this a more generic package.
	scanner := bufio.NewScanner(r)
	lineNum := 0
	records := [][]string{}
	for scanner.Scan() {
		lineNum++
		// TODO: I don't like relying on this Parse() function
		// because if IT breaks then so does this function
		// (i.e. they are coupled). I have some ideas on how
		// to decouple them (in which case this function would
		// become a more generic "map an arbitrary function
		// over each line in a file") but I'm not implementing
		// it in part because I don't think I'll be able to
		// make something as generic as I want because of Go's
		// lack of generics. Also I kind of want to move on
		// with this project and get something submitted so
		// I'll leave it be. By the way, this is a perfect
		// example of why I like functional languages, I feel
		// like they're good about taking an operation which
		// works on one thing and lifting that operation so it
		// works on multiple things.
		record, parseErr := Parse(scanner.Text(), delimiters, numFieldsPerRecord)
		if parseErr != "" {
			parseErrs = append(parseErrs, fmt.Sprintf("%d: %s", lineNum, parseErr))
			continue
		}
		records = append(records, record)
	}
	// TODO: I'm not sure that I like returning a []string when
	// something goes wrong with reading the file. Feels like we
	// should be returning an actual error type, but since we
	// don't do anything different even IF an error were to happen
	// this is fine for now.
	if err := scanner.Err(); err != nil {
		return nil, []string{fmt.Sprintf("unexpected error reading file: %v", err)}
	}
	if len(parseErrs) > 0 {
		return nil, parseErrs
	}
	return records, nil
}
