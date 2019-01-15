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

func parseErrPrefix(lineNum int, msg string) string {
	return fmt.Sprintf("%d: %s", lineNum, msg)
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
	allFields := [][]string{}
	for scanner.Scan() {
		lineNum++
		seps := whichSeparatorsUsedInLine(scanner.Text(), delimiters)
		if len(seps) == 0 {
			parseErrs = append(parseErrs, parseErrPrefix(lineNum, fmt.Sprintf("there is only one field in the record but there should be %d", numFieldsPerRecord)))
			continue
		}
		if len(seps) > 1 {
			sepsStr := fmt.Sprintf("'%c'", seps[0])
			for _, sep := range seps[1:] {
				sepsStr = fmt.Sprintf("%s, '%c'", sepsStr, sep)
			}
			parseErrs = append(parseErrs, parseErrPrefix(lineNum, fmt.Sprintf("there should only be one type of separator in a single line but multiple separators (%s) were specified", sepsStr)))
			continue
		}
		fields := strings.Split(scanner.Text(), string(seps[0]))
		if numFields := len(fields); numFields != numFieldsPerRecord {
			parseErrs = append(parseErrs, parseErrPrefix(lineNum, fmt.Sprintf("there were only %d fields when there should have been %d", numFields, numFieldsPerRecord)))
			continue
		}
		allFields = append(allFields, fields)
	}
	// TODO: I'm not sure that I like returning a []string when
	// something goes wrong with reading the file. Feels like we
	// should be returning an actual error, but since we don't do
	// anything different even IF an error were to happen this is
	// fine for now.
	if err := scanner.Err(); err != nil {
		return nil, []string{fmt.Sprintf("unexpected error reading file: %v", err)}
	}
	if len(parseErrs) > 0 {
		return nil, parseErrs
	}
	return allFields, nil
}
