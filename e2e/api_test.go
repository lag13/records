// +build e2e

package e2e_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	const numRetries = 20
	for i := 0; i < numRetries; i++ {
		if apiIsUp() {
			exitCode := m.Run()
			os.Exit(exitCode)
		}
		log.Print("could not connect to service, taking a short nap before retrying")
		time.Sleep(3 * time.Second)
	}
	log.Fatal("timed out while waiting for service to start")
}

var baseURL = "http://localhost:8080"

func apiIsUp() bool {
	resp, err := http.Get(fmt.Sprintf("%s/healthcheck", baseURL))
	if err != nil {
		fmt.Printf("when checking if the API was up we got an error which was: %v\n", err)
		return false
	}
	resp.Body.Close()
	return http.StatusOK == resp.StatusCode
}

func sendRequest(r *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}
	return resp
}

func newRequest(method string, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseURL, path), body)
	if err != nil {
		panic(err)
	}
	return req
}

func readAll(r io.Reader) []byte {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return b
}

func TestPostFails(t *testing.T) {
	resp := sendRequest(newRequest(http.MethodPost, "/records", strings.NewReader("hey,there buddy")))
	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("when posting invalid data got status code %d, want %d", got, want)
	}
	body := readAll(resp.Body)
	// TODO: Is there even a point to testing that the error
	// message is exactly what we want it to be? I'm not sure that
	// there is. I feel like we should primarily worry about
	// testing things which could be consumed by other programs.
	// These error messages are meant for humans, not machines so
	// if they vary I'm not sure it should matter. Testing that an
	// error message exists seems important but I don't think we
	// have to care overly much about the specific contents.
	if got, want := string(body), `{"errors":["there should only be one type of separator but multiple (',', ' ') were specified"]}`; got != want {
		t.Errorf("when posting invalid data got body %s, want %s", got, want)
	}
}

func TestAPISucceeds(t *testing.T) {
	records := []string{
		"Avatar,Aang,Male,Light-Orange,1760-12-13",
		"MeatAndSarcasmGuy|Sokka|Male|Blue|1845-10-17",
		"SoFullOfHope Katara Female Blue 1846-09-21",
		"BlindBandit,Toph,Female,Green,1846-03-29",
		"Lee|Zuko|Male|Red|1842-07-04",
		"Crazy Azula Female Blood-Red 1842-12-30",
		"Uncle,Iroh,Male,White,1820-08-24",
	}
	for _, record := range records {
		resp := sendRequest(newRequest(http.MethodPost, "/records", strings.NewReader(record)))
		if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Errorf("when posting record %q got status code %d, want %d", record, got, want)
		}
	}
	resp := sendRequest(newRequest(http.MethodGet, "/records/gender", nil))
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("when getting records sorted by gender, got status code %d, wanted %d", got, want)
	}
	body := readAll(resp.Body)
	// TODO: This is just attrocious
	if got, want := string(body), `{"data":[{"last_name":"BlindBandit","first_name":"Toph","gender":"Female","favorite_color":"Green","birthdate":"1846-03-29T00:00:00Z"},{"last_name":"Crazy","first_name":"Azula","gender":"Female","favorite_color":"Blood-Red","birthdate":"1842-12-30T00:00:00Z"},{"last_name":"SoFullOfHope","first_name":"Katara","gender":"Female","favorite_color":"Blue","birthdate":"1846-09-21T00:00:00Z"},{"last_name":"Avatar","first_name":"Aang","gender":"Male","favorite_color":"Light-Orange","birthdate":"1760-12-13T00:00:00Z"},{"last_name":"Lee","first_name":"Zuko","gender":"Male","favorite_color":"Red","birthdate":"1842-07-04T00:00:00Z"},{"last_name":"MeatAndSarcasmGuy","first_name":"Sokka","gender":"Male","favorite_color":"Blue","birthdate":"1845-10-17T00:00:00Z"},{"last_name":"Uncle","first_name":"Iroh","gender":"Male","favorite_color":"White","birthdate":"1820-08-24T00:00:00Z"}]}`; got != want {
		t.Errorf("when getting records sorted by gender, got response body %s, wanted %s", got, want)
	}
}
