package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

func TestExtent_Truncate(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		e := 1*chrono.Hour + 20*chrono.Minute + 5*chrono.Second + 42*chrono.Millisecond + 307*chrono.Microsecond

		if e2 := e.Truncate(200 * chrono.Microsecond); e2 != 4805042200000 {
			t.Fatalf("e.Truncate() = %d, want 4805042200000", e2)
		}
	})

	t.Run("negative", func(t *testing.T) {
		e := -1*chrono.Hour - 20*chrono.Minute - 5*chrono.Second - 42*chrono.Millisecond - 307*chrono.Microsecond

		if e2 := e.Truncate(200 * chrono.Microsecond); e2 != -4805042200000 {
			t.Fatalf("e.Truncate() = %d, want -4805042200000", e2)
		}
	})
}
