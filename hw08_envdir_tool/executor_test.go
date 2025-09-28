package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("bad cmd", func(t *testing.T) {
		exitCode := RunCmd([]string{"ls", "-la", "/notapath"}, nil)
		require.Equal(t, exitCode, 2)
	})
	t.Run("no cmd panic", func(t *testing.T) {
		require.Panics(t, func() {
			emptyArray := make([]string, 0)
			env := Environment{}
			RunCmd(emptyArray, env)
		})
	})
}
