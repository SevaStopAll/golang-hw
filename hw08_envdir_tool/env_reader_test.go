package main

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

type readerTestCase struct {
	name     string
	path     string
	expected Environment
	err      error
}

func TestReadDir(t *testing.T) {
	cases := []readerTestCase{
		{
			name: "simple case based on env testdata",
			path: "./testdata/env",
			expected: Environment{
				"BAR":   EnvValue{Value: "bar", NeedRemove: false},
				"EMPTY": EnvValue{Value: "", NeedRemove: false},
				"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
				"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
				"UNSET": EnvValue{Value: "", NeedRemove: true},
			},
			err: nil,
		},
		{
			name: "case with not existing path",
			path: "/bin/this_path_does_not_exist/",
			expected: Environment{
				"EMPTY_FIRST_LINE": EnvValue{Value: "", NeedRemove: false},
				"ZERO_FILE":        EnvValue{Value: "", NeedRemove: true},
			},
			err: fs.ErrNotExist,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			env, e := ReadDir(tc.path)
			if tc.err != nil {
				require.ErrorIs(t, e, tc.err, "actual error %q but expected %q", e, tc.err)
			} else {
				require.NoError(t, e, "No error expected, but got %q", e)
				require.Equal(t, tc.expected, env)
			}
		})
	}
}
