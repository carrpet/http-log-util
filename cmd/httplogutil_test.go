package cmd

import (
	"fmt"
	"os"
	"testing"
)

var testErrMessage = func(msg, expected, actual string) string {
	return fmt.Sprintf(msg+": "+"expected: %s, actual: %s", expected, actual)
}

func TestArgValidatorInvalidLogFileReturnsDoesNotExistError(t *testing.T) {
	args := []string{"doesntexist.txt"}
	err := argValidator(nil, args)
	if !os.IsNotExist(err) {
		t.Error(
			testErrMessage("Invalid log file did not return an error",
				"argValidator returns PathError",
				"argValidator did not return PathError"))
	}
}

func TestArgValidatorValidFileReturnsNoError(t *testing.T) {
	args := []string{"root.go"}
	err := argValidator(nil, args)
	if err != nil {
		t.Errorf(
			testErrMessage("Existing file returned error: "+err.Error(),
				"argValidator returns nil", "argValidator returns error"))
	}
}

func TestComputeTopHitsReturnsCorrectData(t *testing.T) {
	testdata := [][]string{
		[]string{"10.0.0.2", "-", "apache", "1549573862", "GET /api/user/bleh/h HTTP/1.0", "200", "1234"},
		[]string{"10.0.0.4", "-", "apache", "1549573861", "GET /api/user HTTP/1.0", "200", "1234"},
		[]string{"10.0.0.1", "-", "apache", "1549573862", "GET /api/help HTTP/1.0", "500", "1136"},
		[]string{"10.0.0.4", "-", "apache", "1549573862", "POST /api/help HTTP/1.0", "200", "1234"},
		[]string{"10.0.0.1", "-", "apache", "1549573862", "GET /api/help HTTP/1.0", "200", "1234"},
		[]string{"10.0.0.1", "-", "apache", "1549573862", "GET /report HTTP/1.0", "500", "1194"},
	}
	result := computeTopHits(testdata)
	topSection := result.topHits[0].section
	topHitsCount := result.topHits[0].hits
	if topSection != "/api" {
		t.Errorf(
			testErrMessage("Returned section was incorrect",
				"section should be /api", fmt.Sprintf("section was: %s", topSection)))
	}
	if topHitsCount != "5" {
		t.Errorf(
			testErrMessage("Number of hits was incorrect", "hits should be 5",
				fmt.Sprintf("Hits was %s", topHitsCount)))
	}

}
