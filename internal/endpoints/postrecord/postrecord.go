// Package postrecord parses an POST request which will create a
// record.
package postrecord

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lag13/records/internal/multicsv"
	"github.com/lag13/records/internal/person"
	"github.com/lag13/records/internal/response"
)

// PostRecord parses the incoming request into a person which can then
// be added to the database.
func PostRecord(req *http.Request) (person.Person, response.Structured, error) {
	// TODO: There is repetition in this checking for the correct
	// method and returning an error message if it is not the
	// correct one. One solution would be to use a router which
	// allows you to specify the method when registering the path.
	if req.Method != http.MethodPost {
		return person.Person{}, response.Structured{
			StatusCode: http.StatusBadRequest,
			Errors:     []string{fmt.Sprintf("this endpoint works with a POST request, not a %s", req.Method)},
		}, nil
	}
	r := bufio.NewReader(req.Body)
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		// TODO: If I was being very good I would use
		// pkg/errors to establish a stacktrace at this point
		// in the code so when the error gets logged we know
		// exactly where the failure happened.
		return person.Person{}, response.Structured{
			StatusCode: http.StatusInternalServerError,
			Errors:     []string{"unexpected error"},
		}, err
	}
	line = strings.TrimSpace(line)
	// TODO: I don't like having code, which is unit tested,
	// talking directly to other unit tested code from the same
	// repository because it couples them. But perhaps I'll make
	// an exception with the thought that *this* code, although
	// unit tested, is not going to be consumed by anyone else
	// (except main of course).
	record, parseErr := multicsv.Parse(line, "|, ", 5)
	if parseErr != "" {
		return person.Person{}, response.Structured{
			StatusCode: http.StatusBadRequest,
			Errors:     []string{parseErr},
		}, nil
	}
	p, parseErrs := person.Parse(record)
	if len(parseErrs) > 0 {
		return person.Person{}, response.Structured{
			StatusCode: http.StatusBadRequest,
			Errors:     parseErrs,
		}, nil
	}
	return p, response.Structured{StatusCode: http.StatusOK}, nil
}
