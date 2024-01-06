package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("ls", func(t *testing.T) {
		res := RunCmd([]string{"ls"}, nil)
		require.Zero(t, res)
	})
}
