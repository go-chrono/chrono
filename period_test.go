package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

func TestParseDuration(t *testing.T) {
	t.Run("invalid strings", func(t *testing.T) {
		for _, tt := range []string{
			"P",
			"PT",
		} {
			if _, _, err := chrono.ParseDuration(tt); err == nil {
				t.Fatal("expecting error but got nil")
			}
		}
	})
}
