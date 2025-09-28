package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	path, remaining := cmd[0], make([]string, 0)

	if len(cmd) > 1 {
		remaining = cmd[1:]
	}

	ex := exec.Command(path, remaining...)
	ex.Stdout, ex.Stderr, ex.Stdin = os.Stdout, os.Stderr, os.Stdin

	ex.Env = os.Environ()
	for k, v := range env {
		ex.Env = rmEnv(ex.Env, k)
		if v.NeedRemove {
			continue
		}
		ex.Env = append(ex.Env, fmt.Sprintf("%s=%s", k, v.Value))
	}

	if err := ex.Start(); err != nil {
		return 1
	}

	if err := ex.Wait(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return 1
	}

	return 0
}

func rmEnv(env []string, k string) []string {
	for i, v := range env {
		if strings.Split(v, "=")[0] == k {
			return append(env[:i], env[i+1:]...)
		}
	}

	return env
}
