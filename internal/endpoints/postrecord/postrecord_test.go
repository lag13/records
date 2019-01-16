package postrecord_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/lag13/records/internal/endpoints/postrecord"
	"github.com/lag13/records/internal/response"
)

func errToStr(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprint(err)
}

type mockErrReader struct {
}

func (m mockErrReader) Read([]byte) (int, error) {
	return 0, errors.New("non-nil error")
}

func TestPostRecord(t *testing.T) {
	tests := []struct {
		name     string
		req      *http.Request
		wantResp response.Structured
		errMsg   string
	}{
		{
			name: "invalid http method",
			req:  httptest.NewRequest("GET", "/asdf", nil),
			wantResp: response.Structured{
				StatusCode: 400,
				Errors:     []string{"this endpoint works with a POST request, not a GET"},
			},
			errMsg: "",
		},
		{
			name: "error reading from request",
			req:  httptest.NewRequest("POST", "/asdf", mockErrReader{}),
			wantResp: response.Structured{
				StatusCode: 500,
				Errors:     []string{"unexpected error"},
			},
			errMsg: "non-nil error",
		},
		{
			name: "error parsing line into record",
			req:  httptest.NewRequest("POST", "/asdf", strings.NewReader("hey|there|you")),
			wantResp: response.Structured{
				StatusCode: 400,
				Errors:     []string{"there were 3 fields when there should have been 5"},
			},
			errMsg: "",
		},
		{
			name: "error parsing record to person struct",
			req:  httptest.NewRequest("POST", "/asdf", strings.NewReader("Grey,Gandalf,Male,,1100-04-")),
			wantResp: response.Structured{
				StatusCode: 400,
				Errors:     []string{"favorite color (field 4) must be a non-empty string", "date of birth (field 5) must have the format YYYY-MM-DD"},
			},
			errMsg: "",
		},
		{
			name: "success (and we don't care about other lines)",
			req: httptest.NewRequest("POST", "/asdf", strings.NewReader(`Grey,Gandalf,Male,Rainbow,1100-04-03
this
is
ignored`)),
			wantResp: response.Structured{
				StatusCode: 200,
			},
			errMsg: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := postrecord.PostRecord(test.req)
			if got, want := errToStr(err), test.errMsg; got != want {
				t.Errorf("got error %q, want %q", got, want)
			}
			if got, want := resp, test.wantResp; !reflect.DeepEqual(got, want) {
				t.Errorf("got resp %+v, want %+v", got, want)
			}
			// TODO: Perhaps I should check that the
			// database was updated? Or perhaps I should
			// do that up in main?
		})
	}
}
