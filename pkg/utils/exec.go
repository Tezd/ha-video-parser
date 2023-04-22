package utils

import (
	"io"
	"log"
	"os/exec"
)

const (
	ERROR_EXIT_CODE = 1
)

type (
	ExecResult struct {
		StdOut   string
		StdErr   string
		ExitCode int
	}
)

func Exec(name string, args ...string) ExecResult {
	cmd := exec.Command(name, args...)
	log.Printf("Exec: [%s]", cmd.String())

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return makeErrorResult(err)
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return makeErrorResult(err)
	}

	if err = cmd.Start(); err != nil {
		return makeErrorResult(err)
	}

	outBytes, err := io.ReadAll(stdOut)
	if err != nil {
		return makeErrorResult(err)
	}

	errBytes, err := io.ReadAll(stdErr)
	if err != nil {
		return makeErrorResult(err)
	}

	result := ExecResult{
		StdOut: string(outBytes),
		StdErr: string(errBytes),
	}

	if err = cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = ERROR_EXIT_CODE
		}
	}

	return result
}

func makeErrorResult(err error) ExecResult {
	return ExecResult{
		StdErr:   err.Error(),
		ExitCode: ERROR_EXIT_CODE,
	}
}
