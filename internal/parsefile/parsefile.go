// Package parsefile parses file information so it can be analyzed.
// TODO: Perhaps an informative name for this file could be
// multicharsv if I generalized more.
package parsefile

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// File contains raw data coming from a call to os.Open(). TODO: I
// like that this package gets its data from an io.Reader because it
// makes it more reusable (for example we could read from stdin or
// from the body of a http request). But I do not like the reliance on
// the os package. It's not terrible (you can just ignore the OpenErr
// field when using this package) but it feels a little odd. If I was
// a new person reading this package in isolation I would feel
// confused. In a future iteration I get the feeling we should move
// the logic relating to the os package somewhere else.
type File struct {
	Name    string
	Content io.Reader
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
// analyzed as well as returning errors from parsing. TODO: I would
// prefer returning a []string for the parse errors rather than a
// single string which has newlines for each parse error. What I have
// now works fine for our use case but it feels right to keep the data
// as close as possible to its original form rather than changing it
// around. TODO: I think the final iteration of this function should
// return an error as the third argument in case something related to
// reading the file fails because that is a different kind of "error"
// than invalid structure to a particular line in the file. One is the
// user's fault and the other is something actually went wrong. It's
// sort of like a 4XX vs 5XX status code.
func ParseFile(file File) ([][]string, string) {
	parseErrs := []string{}
	// TODO: I could see us wanting to ignore empty lines but
	// bufio.Scanner does NOT ignore empty lines. Keep this in
	// mind. TODO: Keep in mind that bufio.Scanner has a limited
	// buffer size:
	// https://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go.
	// I don't imagine it would be a problem especially for this
	// practice problem but I wanted to think about it a little
	// more if I tried to make this a more generic package.
	scanner := bufio.NewScanner(file.Content)
	lineNum := 0
	allFields := [][]string{}
	for scanner.Scan() {
		lineNum++
		seps := whichSeparatorsUsedInLine(scanner.Text())
		// TODO: So far with this code all I've been doing is
		// making sure that the file has the expected "shape"
		// without caring about it's semantics similar in
		// spirit to this blog post:
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
