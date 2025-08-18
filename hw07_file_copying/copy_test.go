package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("limit 0 offset 0", func(t *testing.T) {
		err := Copy("testdata\\input.txt",
			"testdata\\test2", 0, 0)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})

	t.Run("offset 0 limit 10", func(t *testing.T) {
		err := Copy("testdata\\input.txt",
			"testdata\\test2", 0, 10)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})

	t.Run("offset 0 limit 1000", func(t *testing.T) {
		err := Copy("testdata\\input.txt",
			"testdata\\test2", 0, 1000)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})

	t.Run("offset 100 limit 1000", func(t *testing.T) {
		err := Copy("testdata\\input.txt",
			"testdata\\test2", 100, 1000)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})

	t.Run("offset 100 limit 1000", func(t *testing.T) {
		err := Copy("testdata\\input.txt",
			"testdata\\test2", 100, 10000)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})

	t.Run("offset 6000 limit 1000", func(t *testing.T) {
		err := Copy("testdata\\input.txt",
			"testdata\\test2", 6000, 1000)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})
}
