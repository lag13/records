package dummy_test

import (
	"testing"

	"github.com/lag13/records/internal/dummy"
)

func TestHelloWorld(t *testing.T) {
	if got, want := dummy.HelloWorld(), "hello world"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
