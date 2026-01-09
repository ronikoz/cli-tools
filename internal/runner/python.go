package runner

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

const defaultPython = "python3"

type RunOptions struct {
	Stream bool
	Python string
}

// RunPython executes a python plugin script and streams output to the console.
func RunPython(scriptPath string, args []string, opts RunOptions) (Result, error) {
	python := opts.Python
	if python == "" {
		python = os.Getenv("CT_PYTHON")
	}
	if python == "" {
		python = defaultPython
	}

	if _, err := os.Stat(scriptPath); err != nil {
		return Result{}, fmt.Errorf("plugin not found: %s", scriptPath)
	}

	cmd := exec.Command(python, append([]string{scriptPath}, args...)...)
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	if opts.Stream {
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	} else {
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
	}
	cmd.Stdin = os.Stdin

	started := time.Now()
	err := cmd.Run()
	finished := time.Now()

	result := Result{
		Command:    python,
		Args:       append([]string{scriptPath}, args...),
		StartedAt:  started,
		FinishedAt: finished,
		DurationMs: finished.Sub(started).Milliseconds(),
		ExitCode:   exitCode(err),
		Stdout:     stdoutBuf.String(),
		Stderr:     stderrBuf.String(),
	}

	if err != nil {
		result.Status = StatusFailed
		result.Error = err.Error()
		return result, errors.New("python runner failed")
	}

	result.Status = StatusSuccess
	return result, nil
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return 1
}


// Signed-off-by: ronikoz
