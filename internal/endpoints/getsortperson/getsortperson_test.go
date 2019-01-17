package getsortperson_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/lag13/records/internal/endpoints/getsortperson"
	"github.com/lag13/records/internal/person"
	"github.com/lag13/records/internal/response"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		req      *http.Request
		sortFn   func(ps []person.Person)
		ps       []person.Person
		wantResp response.Structured
	}{
		{
			name:   "wrong http method",
			req:    httptest.NewRequest("POST", "/asdf", nil),
			sortFn: nil,
			ps:     nil,
			wantResp: response.Structured{
				StatusCode: 400,
				Errors:     []string{"this endpoint works with a GET request, not a POST"},
			},
		},
		{
			name: "do some sorting'ish things!",
			req:  httptest.NewRequest("GET", "/asdf", nil),
			sortFn: func(ps []person.Person) {
				ps[0], ps[1] = ps[1], ps[0]
			},
			ps: []person.Person{
				{LastName: "Bobbo"},
				{LastName: "Vincent"},
			},
			wantResp: response.Structured{
				StatusCode: 200,
				Data: []person.Person{
					{LastName: "Vincent"},
					{LastName: "Bobbo"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := getsortperson.Sort(test.req, test.sortFn, test.ps)
			if got, want := resp, test.wantResp; !reflect.DeepEqual(got, want) {
				t.Errorf("got response %+v, want %+v", got, want)
			}
		})
	}
}
