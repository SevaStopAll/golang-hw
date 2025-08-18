package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("limit 0 offset 0", func(t *testing.T) {
		pathFrom := filepath.Join("testData", "input.txt")
		pathTo := filepath.Join("testdata", "test.txt")
		err := Copy(pathFrom,
			pathTo, 0, 0)
		require.NoError(t, err)
		defer os.Remove(pathTo)
	})

	t.Run("offset 0 limit 10", func(t *testing.T) {
		pathFrom := filepath.Join("testData", "input.txt")
		pathTo := filepath.Join("testdata", "test.txt")
		err := Copy(pathFrom,
			pathTo, 0, 10)
		require.NoError(t, err)
		defer os.Remove(pathTo)
	})

	t.Run("offset 0 limit 1000", func(t *testing.T) {
		pathFrom := filepath.Join("testData", "input.txt")
		pathTo := filepath.Join("testdata", "test.txt")
		err := Copy(pathFrom,
			pathTo, 0, 1000)
		require.NoError(t, err)
		defer os.Remove(pathTo)
	})

	t.Run("offset 100 limit 1000", func(t *testing.T) {
		pathFrom := filepath.Join("testData", "input.txt")
		pathTo := filepath.Join("testdata", "test.txt")
		err := Copy(pathFrom,
			pathTo, 100, 1000)
		require.NoError(t, err)
		defer os.Remove(pathTo)
	})

	t.Run("offset 100 limit 1000", func(t *testing.T) {
		pathFrom := filepath.Join("testData", "input.txt")
		pathTo := filepath.Join("testdata", "test.txt")
		err := Copy(pathFrom,
			pathTo, 100, 10000)
		require.NoError(t, err)
		defer os.Remove(pathTo)
	})

	t.Run("offset 6000 limit 1000", func(t *testing.T) {
		pathFrom := filepath.Join("testData", "input.txt")
		pathTo := filepath.Join("testdata", "test.txt")
		err := Copy(pathFrom,
			pathTo, 6000, 1000)
		require.NoError(t, err)
		defer os.Remove(pathTo)
	})
}
