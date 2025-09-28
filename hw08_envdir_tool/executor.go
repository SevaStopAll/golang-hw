package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envName, envValue := range env {
		if envValue.NeedRemove {
			os.Unsetenv(envName)
		} else {
			os.Setenv(envName, envValue.Value)
		}
	}
	cmdName := cmd[0]  //nolint:gosec
	cmdArgs := cmd[1:] //nolint:gosec
	cmdExec := exec.Command(cmdName, cmdArgs...)
	cmdExec.Env = os.Environ()
	for envName, envValue := range env {
		cmdExec.Env = append(cmdExec.Env, envName+"="+envValue.Value)
	}
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	err := cmdExec.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}
