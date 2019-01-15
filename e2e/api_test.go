// +build e2e

package e2e_test

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

func apiIsUp() bool {
	resp, err := http.Get("http://localhost:8080/healthcheck")
	if err != nil {
		fmt.Printf("when checking if the API was up we got an error which was: %v\n", err)
		return false
	}
	resp.Body.Close()
	return http.StatusOK == resp.StatusCode
}
