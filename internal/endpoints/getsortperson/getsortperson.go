// getsortperson defines a generic handler which will return a
// response of sorted person data.
package getsortperson

import (
	"fmt"
	"net/http"

	"github.com/lag13/records/internal/person"
	"github.com/lag13/records/internal/response"
)

func Sort(req *http.Request, sortFn func(ps []person.Person), ps []person.Person) response.Structured {
	if req.Method != http.MethodGet {
		return response.Structured{
			StatusCode: http.StatusBadRequest,
			Errors:     []string{fmt.Sprintf("this endpoint works with a GET request, not a %s", req.Method)},
		}
	}
	tmp := make([]person.Person, len(ps))
	copy(tmp, ps)
	sortFn(tmp)
	return response.Structured{
		StatusCode: http.StatusOK,
		Data:       tmp,
	}
}
