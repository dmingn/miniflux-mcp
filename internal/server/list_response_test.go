package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	miniflux "miniflux.app/v2/client"
)

type stubClient struct {
	categories miniflux.Categories
	feeds      miniflux.Feeds
}

func (s *stubClient) CategoriesContext(context.Context) (miniflux.Categories, error) {
	return s.categories, nil
}

func (s *stubClient) FeedsContext(context.Context) (miniflux.Feeds, error) {
	return s.feeds, nil
}

func (s *stubClient) CategoryFeedsContext(context.Context, int64) (miniflux.Feeds, error) {
	return s.feeds, nil
}

func (s *stubClient) FeedContext(context.Context, int64) (*miniflux.Feed, error) {
	return nil, nil
}

func (s *stubClient) EntriesContext(context.Context, *miniflux.Filter) (*miniflux.EntryResultSet, error) {
	return nil, nil
}

func (s *stubClient) EntryContext(context.Context, int64) (*miniflux.Entry, error) {
	return nil, nil
}

func (s *stubClient) FetchCountersContext(context.Context) (*miniflux.FeedCounters, error) {
	return nil, nil
}

func TestListCategories_returnsObjectStructuredContent(t *testing.T) {
	t.Parallel()

	// Arrange
	h := &handlers{client: &stubClient{
		categories: miniflux.Categories{{ID: 1, Title: "News"}},
	}}

	// Act
	_, structured, err := h.listCategories(context.Background(), &mcp.CallToolRequest{}, struct{}{})

	// Assert
	if err != nil {
		t.Fatalf("listCategories: %v", err)
	}
	assertJSONObjectWithArrayKey(t, structured, "categories")
}

func TestListFeeds_returnsObjectStructuredContent(t *testing.T) {
	t.Parallel()

	// Arrange
	h := &handlers{client: &stubClient{
		feeds: miniflux.Feeds{{ID: 2, Title: "Example"}},
	}}

	// Act
	_, structured, err := h.listFeeds(context.Background(), &mcp.CallToolRequest{}, listFeedsInput{})

	// Assert
	if err != nil {
		t.Fatalf("listFeeds: %v", err)
	}
	assertJSONObjectWithArrayKey(t, structured, "feeds")
}

func assertJSONObjectWithArrayKey(t *testing.T, value any, key string) {
	t.Helper()

	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if len(raw) == 0 || raw[0] != '{' {
		t.Fatalf("structuredContent must be a JSON object, got %s", raw)
	}

	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal object: %v", err)
	}
	items, ok := decoded[key]
	if !ok {
		t.Fatalf("missing %q key in %s", key, raw)
	}
	if len(items) == 0 || items[0] != '[' {
		t.Fatalf("%q must be a JSON array, got %s", key, items)
	}
}
