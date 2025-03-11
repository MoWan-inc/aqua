package util

import (
	"encoding/json"
	"fmt"
	"os"
)

type errorMessage struct {
	Msg      string `json:"msg"`
	ExitCode int    `json:"exit_code"`
	Error    string `json:"error"`
}

func ErrorMessage(exitCode int, err error) string {
	msg := errorMessage{
		Msg:      "component exit abnormally",
		ExitCode: exitCode,
		Error:    err.Error(),
	}
	msgData, _ := json.Marshal(msg)
	return string(msgData)
}

func ExitPrintFatalError(err error) {
	var exitCode int
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, ErrorMessage(-1, err))
		exitCode = -1
	}
	os.Exit(exitCode)
}
