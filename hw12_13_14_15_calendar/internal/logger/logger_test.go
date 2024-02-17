package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("test log", func(t *testing.T) {
		l, _ := New("Info")
		l.Info("Info")
		l.Warn("Warn")
		l.Error("Error")
	})
}
