package server

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	miniflux "miniflux.app/v2/client"

	"github.com/dmingn/miniflux-mcp/internal/preview"
	"github.com/dmingn/miniflux-mcp/internal/version"
)

var _ Client = (*miniflux.Client)(nil)

type Client interface {
	CategoriesContext(ctx context.Context) (miniflux.Categories, error)
	FeedsContext(ctx context.Context) (miniflux.Feeds, error)
	CategoryFeedsContext(ctx context.Context, categoryID int64) (miniflux.Feeds, error)
	FeedContext(ctx context.Context, feedID int64) (*miniflux.Feed, error)
	EntriesContext(ctx context.Context, filter *miniflux.Filter) (*miniflux.EntryResultSet, error)
	EntryContext(ctx context.Context, entryID int64) (*miniflux.Entry, error)
	FetchCountersContext(ctx context.Context) (*miniflux.FeedCounters, error)
}

func New(client Client) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "miniflux-mcp",
		Version: version.String(),
	}, nil)

	h := &handlers{client: client}
	h.register(s)
	return s
}

type handlers struct {
	client Client
}

type listFeedsInput struct {
	CategoryID *int64 `json:"category_id,omitempty" jsonschema:"optional category ID to filter feeds"`
}

type getFeedInput struct {
	FeedID int64 `json:"feed_id" jsonschema:"feed ID"`
}

type getEntriesInput struct {
	Status     string `json:"status,omitempty" jsonschema:"entry status filter: unread, read, or removed; default unread"`
	FeedID     *int64 `json:"feed_id,omitempty" jsonschema:"filter by feed ID"`
	CategoryID *int64 `json:"category_id,omitempty" jsonschema:"filter by category ID"`
	Search     string `json:"search,omitempty" jsonschema:"full-text search query"`
	Limit      *int   `json:"limit,omitempty" jsonschema:"maximum number of entries to return; default 20"`
	Offset     *int   `json:"offset,omitempty" jsonschema:"pagination offset; default 0"`
	Order      string `json:"order,omitempty" jsonschema:"sort field: id, status, published_at, category_title, category_id"`
	Direction  string `json:"direction,omitempty" jsonschema:"sort direction: asc or desc"`
	Starred    *bool  `json:"starred,omitempty" jsonschema:"if set, filter to starred or non-starred entries"`
}

type getEntryInput struct {
	EntryID int64 `json:"entry_id" jsonschema:"entry ID"`
}

func (h *handlers) register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_categories",
		Description: "List Miniflux categories",
	}, h.listCategories)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_feeds",
		Description: "List Miniflux feeds. Optionally filter by category_id.",
	}, h.listFeeds)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_feed",
		Description: "Get a single Miniflux feed by ID",
	}, h.getFeed)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_entries",
		Description: "List Miniflux entries. Defaults: status=unread, limit=20. Does not include full HTML content; each entry has content_preview instead.",
	}, h.getEntries)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_entry",
		Description: "Get a single Miniflux entry by ID, including full HTML content",
	}, h.getEntry)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_counters",
		Description: "Get read and unread entry counters per feed",
	}, h.getCounters)
}

func (h *handlers) listCategories(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	categories, err := h.client.CategoriesContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	return nil, categories, nil
}

func (h *handlers) listFeeds(ctx context.Context, _ *mcp.CallToolRequest, input listFeedsInput) (*mcp.CallToolResult, any, error) {
	if input.CategoryID != nil {
		if *input.CategoryID <= 0 {
			return nil, nil, fmt.Errorf("category_id must be a positive integer")
		}
		feeds, err := h.client.CategoryFeedsContext(ctx, *input.CategoryID)
		if err != nil {
			return nil, nil, err
		}
		return nil, feeds, nil
	}

	feeds, err := h.client.FeedsContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	return nil, feeds, nil
}

func (h *handlers) getFeed(ctx context.Context, _ *mcp.CallToolRequest, input getFeedInput) (*mcp.CallToolResult, any, error) {
	if input.FeedID <= 0 {
		return nil, nil, fmt.Errorf("feed_id must be a positive integer")
	}

	feed, err := h.client.FeedContext(ctx, input.FeedID)
	if err != nil {
		return nil, nil, err
	}
	return nil, feed, nil
}

func (h *handlers) getEntries(ctx context.Context, _ *mcp.CallToolRequest, input getEntriesInput) (*mcp.CallToolResult, any, error) {
	filter := &miniflux.Filter{
		Status: miniflux.EntryStatusUnread,
		Limit:  20,
	}
	if input.Status != "" {
		filter.Status = input.Status
	}
	if input.Limit != nil {
		if *input.Limit < 0 {
			return nil, nil, fmt.Errorf("limit must be >= 0")
		}
		filter.Limit = *input.Limit
	}
	if input.Offset != nil {
		if *input.Offset < 0 {
			return nil, nil, fmt.Errorf("offset must be >= 0")
		}
		filter.Offset = *input.Offset
	}
	if input.FeedID != nil {
		if *input.FeedID <= 0 {
			return nil, nil, fmt.Errorf("feed_id must be a positive integer")
		}
		filter.FeedID = *input.FeedID
	}
	if input.CategoryID != nil {
		if *input.CategoryID <= 0 {
			return nil, nil, fmt.Errorf("category_id must be a positive integer")
		}
		filter.CategoryID = *input.CategoryID
	}
	if input.Search != "" {
		filter.Search = input.Search
	}
	if input.Order != "" {
		filter.Order = input.Order
	}
	if input.Direction != "" {
		filter.Direction = input.Direction
	}
	if input.Starred != nil {
		if *input.Starred {
			filter.Starred = miniflux.FilterOnlyStarred
		} else {
			filter.Starred = miniflux.FilterNotStarred
		}
	}

	result, err := h.client.EntriesContext(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	return nil, entryListResponse{
		Total:   result.Total,
		Entries: toEntrySummaries(result.Entries),
	}, nil
}

func (h *handlers) getEntry(ctx context.Context, _ *mcp.CallToolRequest, input getEntryInput) (*mcp.CallToolResult, any, error) {
	if input.EntryID <= 0 {
		return nil, nil, fmt.Errorf("entry_id must be a positive integer")
	}

	entry, err := h.client.EntryContext(ctx, input.EntryID)
	if err != nil {
		return nil, nil, err
	}
	return nil, entry, nil
}

func (h *handlers) getCounters(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	counters, err := h.client.FetchCountersContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	return nil, counters, nil
}

type entryListResponse struct {
	Total   int            `json:"total"`
	Entries []entrySummary `json:"entries"`
}

type entrySummary struct {
	ID             int64               `json:"id"`
	UserID         int64               `json:"user_id"`
	FeedID         int64               `json:"feed_id"`
	Status         string              `json:"status"`
	Hash           string              `json:"hash"`
	Title          string              `json:"title"`
	URL            string              `json:"url"`
	CommentsURL    string              `json:"comments_url"`
	PublishedAt    time.Time           `json:"published_at"`
	CreatedAt      time.Time           `json:"created_at"`
	ChangedAt      time.Time           `json:"changed_at"`
	Author         string              `json:"author"`
	Language       string              `json:"language"`
	ShareCode      string              `json:"share_code"`
	Starred        bool                `json:"starred"`
	ReadingTime    int                 `json:"reading_time"`
	Tags           []string            `json:"tags"`
	ContentPreview string              `json:"content_preview"`
	Enclosures     miniflux.Enclosures `json:"enclosures,omitempty"`
}

func toEntrySummaries(entries miniflux.Entries) []entrySummary {
	out := make([]entrySummary, 0, len(entries))
	for _, e := range entries {
		if e == nil {
			continue
		}
		out = append(out, entrySummary{
			ID:             e.ID,
			UserID:         e.UserID,
			FeedID:         e.FeedID,
			Status:         e.Status,
			Hash:           e.Hash,
			Title:          e.Title,
			URL:            e.URL,
			CommentsURL:    e.CommentsURL,
			PublishedAt:    e.Date,
			CreatedAt:      e.CreatedAt,
			ChangedAt:      e.ChangedAt,
			Author:         e.Author,
			Language:       e.Language,
			ShareCode:      e.ShareCode,
			Starred:        e.Starred,
			ReadingTime:    e.ReadingTime,
			Tags:           e.Tags,
			ContentPreview: preview.ContentPreview(e.Content, 100),
			Enclosures:     e.Enclosures,
		})
	}
	return out
}
