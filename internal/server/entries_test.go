package server

import (
	"testing"

	miniflux "miniflux.app/v2/client"
)

func TestToEntrySummaries_replacesContentWithPreview(t *testing.T) {
	t.Parallel()

	// Arrange
	entries := miniflux.Entries{
		{
			ID:      1,
			Title:   "Hello",
			Content: "<p>Hello <b>world</b></p>",
		},
	}

	// Act
	summaries := toEntrySummaries(entries)

	// Assert
	if len(summaries) != 1 {
		t.Fatalf("got %d summaries, want 1", len(summaries))
	}
	if summaries[0].ContentPreview != "Hello world" {
		t.Fatalf("preview = %q", summaries[0].ContentPreview)
	}
}
