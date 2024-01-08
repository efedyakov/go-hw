package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	dir := "testdata/env"

	t.Run("incorrect dir", func(t *testing.T) {
		incorrectDir := filepath.Join(dir, "incorrect_dir")
		_, err := ReadDir(incorrectDir)
		require.Error(t, err)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("not dir", func(t *testing.T) {
		filename := filepath.Join(dir, "file.txt")
		_, _ = os.Create(filename)
		defer os.Remove(filename)

		_, err := ReadDir(filename)
		require.Error(t, err)
	})

	t.Run("empty dir", func(t *testing.T) {
		emptydir := filepath.Join(dir, "emptydir")
		_ = os.Mkdir(emptydir, os.ModePerm)
		defer os.Remove(emptydir)

		env, err := ReadDir(emptydir)
		require.NoError(t, err)
		require.Empty(t, env)
	})

	t.Run("dir", func(t *testing.T) {
		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, true, env["UNSET"].NeedRemove)
		require.Equal(t, "", env["EMPTY"].Value)
		require.Equal(t, "foo\nwith new line", env["FOO"].Value)
		require.NotEmpty(t, env)
	})
}
