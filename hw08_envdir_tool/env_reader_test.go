package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testReadDir      = "testdata1"
	testReadFileName = "ENV=1"
	testFile         *os.File
	testReadExpEnv   = "ENV1"
)

func TestReadDir(t *testing.T) {
	t.Run("Test env reader with '=' char", func(t *testing.T) {
		err := os.Mkdir(testReadDir, 0o777)
		if err != nil {
			panic(err)
		}
		testFile, err = os.Create(testReadDir + "/" + testReadFileName)
		if err != nil {
			os.Remove(testReadDir)
			panic(err)
		}
		defer testReadDirClean()
		testFile.Close()
		envMap, err := ReadDir(testReadDir)
		require.NoError(t, err)
		envValue, ok := envMap[testReadExpEnv]
		require.True(t, ok)
		require.Equal(t, "", envValue.Value)
		require.Equal(t, true, envValue.NeedRemove)
	})
}

func testReadDirClean() {
	os.RemoveAll(testReadDir)
}
