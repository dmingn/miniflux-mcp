package version

import "testing"

func TestString_returnsNonEmpty(t *testing.T) {
	t.Parallel()

	// Act
	got := String()

	// Assert
	if got == "" {
		t.Fatal("expected non-empty version")
	}
}
