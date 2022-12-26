package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

func TestDivideAndRoundInt(t *testing.T) {
	t.Run("round up", func(t *testing.T) {
		if rounded := chrono.DivideAndRoundIntFunc(123456, 1000); rounded != 123 {
			t.Errorf("RoundInt() = %d, want %d", rounded, 123)
		}
	})

	t.Run("round down", func(t *testing.T) {
		if rounded := chrono.DivideAndRoundIntFunc(123456789, 1000); rounded != 123457 {
			t.Errorf("RoundInt() = %d, want %d", rounded, 123457)
		}
	})
}
