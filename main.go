package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lag13/records/internal/parsefile"
)

func main() {
	parseErrs := []string{}
	{ // parse each file
		for _, fileName := range os.Args[1:] {
			fh, err := os.Open(fileName)
			parseErr := parsefile.ParseFile(parsefile.File{
				Name:    fileName,
				Content: fh,
				OpenErr: err,
			})
			if parseErr != "" {
				parseErrs = append(parseErrs, parseErr)
			}
			if fh != nil {
				fh.Close()
			}
		}
	}
	if len(parseErrs) > 0 {
		fmt.Println(strings.Join(parseErrs, "\n"))
		return
	}
}
