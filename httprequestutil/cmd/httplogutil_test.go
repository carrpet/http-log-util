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
				"argValidator returns nil", "argValidator returns err"))
	}
}
