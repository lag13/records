package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lag13/records/internal/parsefile"
	"github.com/lag13/records/internal/person"
)

func main() {
	persons := []person.Person{}
	parseErrs := []string{}
	// parse each file. TODO: This is getting QUITE messy and
	// feels a little too complicated for being in main. I'm
	// thinking that I should break this up into separate
	// structural parsing and semantic parsing portions. Or
	// perhaps I'll try to move the code around a little (maybe
	// make a new unit that can be unit tested?). But! tests pass
	// so at least we can refactor with a little confidence.
	{
		for _, fileName := range os.Args[1:] {
			fh, err := os.Open(fileName)
			defer func(fh *os.File) {
				if fh == nil {
					return
				}
				_ = fh.Close()
			}(fh)
			lines, parseErr := parsefile.ParseFile(parsefile.File{
				Name:    fileName,
				Content: fh,
				OpenErr: err,
			})
			if parseErr != "" {
				parseErrs = append(parseErrs, parseErr)
				continue
			}
			for i, line := range lines {
				p, semParseErrs := person.Parse(line)
				annotatedSemParseErrs := []string{}
				for _, semParseErr := range semParseErrs {
					annotatedSemParseErrs = append(annotatedSemParseErrs, fmt.Sprintf("%s:%d: %s", fileName, i+1, semParseErr))
				}
				parseErrs = append(parseErrs, annotatedSemParseErrs...)
				persons = append(persons, p)
			}
		}
	}
	if len(parseErrs) > 0 {
		fmt.Println("Invalid Input:")
		fmt.Println(strings.Join(parseErrs, "\n"))
		return
	}
	person.SortByGenderThenLastName(persons)
	for _, p := range persons {
		fmt.Println(person.Marshal(p))
	}
}
