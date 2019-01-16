// postrecord adds a record to the database
package postrecord

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lag13/records/internal/db"
	"github.com/lag13/records/internal/multicsv"
	"github.com/lag13/records/internal/person"
	"github.com/lag13/records/internal/response"
)

// PostRecord adds a record to the database.
func PostRecord(req *http.Request) (response.Structured, error) {
	if req.Method != http.MethodPost {
		return response.Structured{
			StatusCode: http.StatusBadRequest,
			Errors:     []string{fmt.Sprintf("this endpoint works with a POST request, not a %s", req.Method)},
		}, nil
	}
	r := bufio.NewReader(req.Body)
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		// TODO: If I was being very good this might be where
		// I use pkg/errors to establish a stacktrace at this
		// point in the code so when the error gets logged we
		// know exactly where the failure happened.
		return response.Structured{
			StatusCode: http.StatusInternalServerError,
			Errors:     []string{"unexpected error"},
		}, err
	}
	line = strings.TrimSpace(line)
	// TODO: Again, I don't like having other bits of code which
	// are unit tested talking directly to other unit tested code
	// in the same repository because it couples them. But perhaps
	// I'll make an exception with the thought that *this* code,
	// although unit tested, is not going to be consumed by anyone
	// else (except main of course)
	record, parseErr := multicsv.Parse(line, "|, ", 5)
	if parseErr != "" {
		return response.Structured{
			StatusCode: http.StatusBadRequest,
			Errors:     []string{parseErr},
		}, nil
	}
	p, parseErrs := person.Parse(record)
	if len(parseErrs) > 0 {
		return response.Structured{
			StatusCode: http.StatusBadRequest,
			Errors:     parseErrs,
		}, nil
	}
	// TODO: Mutex of just return the person and let someone else
	// deal with it.
	db.Persons = append(db.Persons, p)
	return response.Structured{StatusCode: http.StatusOK}, nil
}
