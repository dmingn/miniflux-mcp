package preview

import (
	"strings"
	"testing"
)

func TestContentPreview_stripsHTML(t *testing.T) {
	t.Parallel()

	// Arrange
	raw := "<p>Hello <b>world</b> &amp; friends</p>"
	want := "Hello world & friends"

	// Act
	got := ContentPreview(raw, 100)

	// Assert
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestContentPreview_keepsTextWhenGreaterThanAppearsOutsideTags(t *testing.T) {
	t.Parallel()

	// Arrange
	// A naive tag regex (<[^>]*>) would cut at the first ">" inside the comment
	// or attribute and leave junk like " this -->" / "b\">ok".
	raw := `<!-- ignore > this --><p title="a>b">ok</p>`
	want := "ok"

	// Act
	got := ContentPreview(raw, 100)

	// Assert
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestContentPreview_runeAwareTruncation(t *testing.T) {
	t.Parallel()

	// Arrange
	raw := "<div>あいうえおかきくけこ</div>"
	want := "あいうえお"

	// Act
	got := ContentPreview(raw, 5)

	// Assert
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestContentPreview_defaultMaxRunes(t *testing.T) {
	t.Parallel()

	// Arrange
	raw := strings.Repeat("あ", 150)

	// Act
	got := ContentPreview(raw, 0)

	// Assert
	if len([]rune(got)) != 100 {
		t.Fatalf("got %d runes, want 100", len([]rune(got)))
	}
}
