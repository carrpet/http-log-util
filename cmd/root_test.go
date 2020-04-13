package cmd

import (
	"os"
	"testing"
)

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
