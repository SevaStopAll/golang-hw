package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("limit 0 offset 0", func(t *testing.T) {
		fmt.Println(os.Getwd())
		pathFrom := filepath.Join("testdata", "input.txt")
		pathTo := filepath.Join("testdata", "outtest.txt")
		_, errorCreating := os.CreateTemp("", pathTo)
		if errorCreating != nil {
			return
		}
		err := Copy(pathFrom,
			pathTo, 0, 0)

		equalFile, err2 := os.Stat(filepath.Join("testdata", "out_offset0_limit0.txt"))
		if err2 != nil {
			t.Fatal(err2)
		}
		writtenFile, err2 := os.Stat(pathTo)
		if err2 != nil {
			t.Fatal(err2)
		}

		require.NoError(t, err)
		require.Equal(t, equalFile.Size(), writtenFile.Size())
		defer os.Remove(pathTo)
	})

	t.Run("offset 0 limit 10", func(t *testing.T) {
		pathFrom := filepath.Join("testdata", "input.txt")
		pathTo := filepath.Join("testdata", "outtest.txt")
		_, errorCreating := os.CreateTemp("", pathTo)
		if errorCreating != nil {
			return
		}
		err := Copy(pathFrom,
			pathTo, 0, 10)

		equalFile, err2 := os.Stat(filepath.Join("testdata", "out_offset0_limit10.txt"))
		if err2 != nil {
			t.Fatal(err2)
		}
		writtenFile, err2 := os.Stat(pathTo)
		if err2 != nil {
			t.Fatal(err2)
		}

		require.NoError(t, err)
		require.Equal(t, equalFile.Size(), writtenFile.Size())
		defer os.Remove(pathTo)
	})

	t.Run("offset 0 limit 1000", func(t *testing.T) {
		pathFrom := filepath.Join("testdata", "input.txt")
		pathTo := filepath.Join("testdata", "outtest.txt")

		_, errorCreating := os.CreateTemp("", pathTo)
		if errorCreating != nil {
			return
		}

		err := Copy(pathFrom,
			pathTo, 0, 1000)

		equalFile, err2 := os.Stat(filepath.Join("testdata", "out_offset0_limit1000.txt"))
		if err2 != nil {
			t.Fatal(err2)
		}
		writtenFile, err2 := os.Stat(pathTo)
		if err2 != nil {
			t.Fatal(err2)
		}

		require.NoError(t, err)
		require.Equal(t, equalFile.Size(), writtenFile.Size())
		defer os.Remove(pathTo)
	})

	t.Run("offset 100 limit 1000", func(t *testing.T) {
		pathFrom := filepath.Join("testdata", "input.txt")
		pathTo := filepath.Join("testdata", "outtest.txt")

		_, errorCreating := os.CreateTemp("", pathTo)
		if errorCreating != nil {
			return
		}

		err := Copy(pathFrom,
			pathTo, 100, 1000)
		equalFile, err2 := os.Stat(filepath.Join("testdata", "out_offset100_limit1000.txt"))

		if err2 != nil {
			t.Fatal(err2)
		}
		writtenFile, err2 := os.Stat(pathTo)
		if err2 != nil {
			t.Fatal(err2)
		}

		require.NoError(t, err)
		require.Equal(t, equalFile.Size(), writtenFile.Size())
		defer os.Remove(pathTo)
	})

	t.Run("offset 6000 limit 1000", func(t *testing.T) {
		pathFrom := filepath.Join("testdata", "input.txt")
		pathTo := filepath.Join("testdata", "outtest.txt")

		_, errorCreating := os.CreateTemp("", pathTo)
		if errorCreating != nil {
			return
		}

		err := Copy(pathFrom,
			pathTo, 6000, 1000)

		equalFile, err2 := os.Stat(filepath.Join("testdata", "out_offset6000_limit1000.txt"))
		if err2 != nil {
			t.Fatal(err2)
		}
		writtenFile, err2 := os.Stat(pathTo)
		if err2 != nil {
			t.Fatal(err2)
		}

		require.NoError(t, err)
		require.Equal(t, equalFile.Size(), writtenFile.Size())
		defer os.Remove(pathTo)
	})
}
