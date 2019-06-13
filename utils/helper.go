package utils

import (
	"fmt"
	"runtime"
	"strings"
)

// CatchPanic is used to catch any Panic and log exceptions to Stdout. It will also write the stack trace
func CatchPanic(err *error, sessionID string, functionName string) {
	if r := recover(); r != nil {
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

const (
	// MainGoRoutine is just a label for logging.
	MainGoRoutine = "main"
)

func CheckStringIsBlank(s string) bool {
	if len(s) <= 0 || s == "" || strings.Count(s, " ") == len(s) {
		return true
	}
	return false
}
