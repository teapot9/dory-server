package helpers

import (
	"testing"
)

func logReq(t *testing.T, payload any, response []byte) {
	t.Logf("sent HTTP request with payload: %v", payload)
	t.Logf("server response was: %s", string(response))
}

// AssertHTTPResponse assert that the HTTP request has no errors and
// had the expected response code
func AssertHTTPResponse(
	t *testing.T, msg string, payload any, response []byte,
	err error, codeExp, codeReal int,
) {
	if err != nil {
		logReq(t, payload, response)
		t.Fatalf("%s: error: %v", msg, err)
	}
	if codeReal != codeExp {
		logReq(t, payload, response)
		t.Fatalf("%s: expected HTTP %d but got %d", msg, codeExp, codeReal)
	}
}
