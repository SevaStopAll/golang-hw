package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("limit 0 offset 0", func(t *testing.T) {
		err := Copy("D:\\GolandProjects\\golang-hw\\hw07_file_copying\\testdata\\out_offset0_limit0.txt",
			"testdata\\test2", 0, 0)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})

	t.Run("offset 10 limit 0", func(t *testing.T) {
		err := Copy("D:\\GolandProjects\\golang-hw\\hw07_file_copying\\testdata\\out_offset0_limit10.txt",
			"testdata\\test2", 0, 1024)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})

	t.Run("offset 0 limit 1000", func(t *testing.T) {
		err := Copy("D:\\GolandProjects\\golang-hw\\hw07_file_copying\\testdata\\out_offset0_limit1000.txt",
			"testdata\\test2", 0, 1024)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})

	t.Run("offset 100 limit 1000", func(t *testing.T) {
		err := Copy("D:\\GolandProjects\\golang-hw\\hw07_file_copying\\testdata\\out_offset100_limit1000.txt",
			"testdata\\test2", 0, 1024)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})

	t.Run("offset 6000 limit 1000", func(t *testing.T) {
		err := Copy("D:\\GolandProjects\\golang-hw\\hw07_file_copying\\testdata\\out_offset6000_limit1000.txt",
			"testdata\\test2", 0, 1024)
		require.NoError(t, err)
		defer os.Remove("testdata\\test2")
	})
}
