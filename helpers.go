package tavily

import (
	"context"
	"fmt"
)

// SearchSimple performs a basic search with minimal configuration.
// It's a convenience method for quick searches without configuring options.
func (c *Client) SearchSimple(ctx context.Context, query string) (*SearchResponse, error) {
	return c.Search(ctx, query, nil)
}

// SearchWithAnswer performs a search and requests an AI-generated answer.
func (c *Client) SearchWithAnswer(ctx context.Context, query string) (*SearchResponse, error) {
	opts := &SearchOptions{
		IncludeAnswer: true,
		MaxResults:    10,
	}
	return c.Search(ctx, query, opts)
}

// SearchNews performs a news-focused search with recent results.
func (c *Client) SearchNews(ctx context.Context, query string, days int) (*SearchResponse, error) {
	opts := &SearchOptions{
		Topic:         string(TopicNews),
		SearchDepth:   string(SearchDepthAdvanced),
		Days:          days,
		MaxResults:    15,
		IncludeAnswer: true,
	}
	return c.Search(ctx, query, opts)
}

// ExtractSimple extracts content from a single URL with default settings.
func (c *Client) ExtractSimple(ctx context.Context, url string) (*ExtractResponse, error) {
	return c.Extract(ctx, []string{url}, nil)
}

// ExtractWithImages extracts content and images from URLs.
func (c *Client) ExtractWithImages(ctx context.Context, urls []string) (*ExtractResponse, error) {
	opts := &ExtractOptions{
		IncludeImages: BoolPtr(true),
		Format:        string(FormatMarkdown),
		ExtractDepth:  string(SearchDepthAdvanced),
	}
	return c.Extract(ctx, urls, opts)
}

// CrawlDocumentation crawls a website focusing on documentation pages.
func (c *Client) CrawlDocumentation(ctx context.Context, url string, maxPages int) (*CrawlResponse, error) {
	opts := &CrawlOptions{
		MaxDepth:      3,
		Limit:         maxPages,
		Categories:    []CrawlCategory{CategoryDocumentation, CategoryDeveloper},
		SelectPaths:   []string{"/docs/*", "/api/*", "/guide/*", "/tutorial/*"},
		Format:        string(FormatMarkdown),
		AllowExternal: BoolPtr(false),
	}
	return c.Crawl(ctx, url, opts)
}

// MapSite provides a quick way to map a website structure.
func (c *Client) MapSite(ctx context.Context, url string) (*MapResponse, error) {
	opts := &MapOptions{
		MaxDepth: 2,
		Limit:    100,
	}
	return c.Map(ctx, url, opts)
}

// GetSearchContext returns search results formatted as context for AI applications.
// This is useful for RAG (Retrieval-Augmented Generation) workflows.
func (c *Client) GetSearchContext(ctx context.Context, query string, maxTokens int) (string, error) {
	opts := &SearchOptions{
		SearchDepth:       string(SearchDepthAdvanced),
		MaxResults:        5,
		IncludeRawContent: string(FormatText),
		MaxTokens:         maxTokens,
	}

	result, err := c.Search(ctx, query, opts)
	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}

	context := fmt.Sprintf("Search query: %s\n\n", query)
	for i, r := range result.Results {
		context += fmt.Sprintf("Source %d: %s\nURL: %s\nContent: %s\n\n",
			i+1, r.Title, r.URL, r.Content)
	}

	return context, nil
}

// BoolPtr is a helper function to get a pointer to a boolean value.
// This is useful for optional boolean fields in API requests.
func BoolPtr(b bool) *bool {
	return &b
}

// GetVersionInfo returns version information about the client.
func GetVersionInfo() map[string]string {
	return map[string]string{
		"client_name":    "go-tavily",
		"client_version": "1.0.0",
		"go_version":     "1.24+",
		"api_version":    "v1",
	}
}
