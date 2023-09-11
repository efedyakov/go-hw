package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	sfile := "testdata/input.txt"
	dfile := "testdata/input2.txt"

	t.Run("simple copy", func(t *testing.T) {
		err := Copy(sfile, dfile, 0, 0)
		require.Empty(t, err)
		require.FileExists(t, dfile)

		cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
		equal, _ := cmp.CompareFile(sfile, dfile)
		require.EqualValues(t, true, equal)
	})

	t.Run("copy 100", func(t *testing.T) {
		err := Copy(sfile, dfile, 0, 100)
		require.Empty(t, err)
		require.FileExists(t, dfile)

		cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
		equal, _ := cmp.CompareFile(sfile, dfile)
		require.EqualValues(t, false, equal)
	})

	t.Run("copy 100 offset 100", func(t *testing.T) {
		err := Copy(sfile, dfile, 100, 200)
		require.Empty(t, err)
		require.FileExists(t, dfile)

		cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
		equal, _ := cmp.CompareFile(sfile, dfile)
		require.EqualValues(t, false, equal)
	})
}
