// Package tavily provides a modern Go client for the Tavily AI-powered search and web content extraction API.
//
// The client supports all four main Tavily API operations:
// - Search: Web search with intelligent results aggregation
// - Extract: Content extraction from specific URLs
// - Crawl: Intelligent website crawling and content mapping
// - Map: Website structure discovery and mapping
//
// Usage:
//
//	client := tavily.New("tvly-your-api-key")
//	result, err := client.Search(ctx, "Go programming language", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Found %d results\n", len(result.Results))
package tavily

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	DefaultBaseURL     = "https://api.tavily.com"
	DefaultTimeout     = 60 * time.Second
	DefaultMaxResults  = 5
	DefaultSearchDepth = "basic"
	DefaultTopic       = "general"
	DefaultFormat      = "text"
	ClientSource       = "go-tavily"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	headers    map[string]string
}

type Options struct {
	BaseURL    string
	HTTPClient *http.Client
	Timeout    time.Duration
}

// New creates a new Tavily API client with the provided API key.
// If apiKey is empty, it attempts to read from TAVILY_API_KEY environment variable.
func New(apiKey string, opts *Options) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("TAVILY_API_KEY")
	}

	if opts == nil {
		opts = &Options{}
	}

	baseURL := opts.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: timeout,
		}
	}

	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: httpClient,
		headers: map[string]string{
			"Content-Type":    "application/json",
			"Authorization":   "Bearer " + apiKey,
			"X-Client-Source": ClientSource,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, endpoint string, requestBody any, responseBody any) error {
	if c.apiKey == "" {
		return &APIError{
			StatusCode: 401,
			Message:    "missing API key - provide via parameter or TAVILY_API_KEY environment variable",
		}
	}

	var body io.Reader
	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return parseAPIError(resp.StatusCode, respData)
	}

	if responseBody != nil {
		if err := json.Unmarshal(respData, responseBody); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func parseAPIError(statusCode int, respData []byte) error {
	var errorResp struct {
		Detail struct {
			Error string `json:"error"`
		} `json:"detail"`
	}

	message := "unknown error"
	if json.Unmarshal(respData, &errorResp) == nil && errorResp.Detail.Error != "" {
		message = errorResp.Detail.Error
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// Search performs an intelligent web search with advanced filtering and content aggregation.
func (c *Client) Search(ctx context.Context, query string, opts *SearchOptions) (*SearchResponse, error) {
	if opts == nil {
		opts = &SearchOptions{}
	}

	req := &SearchRequest{
		Query:                    query,
		SearchDepth:              defaultString(opts.SearchDepth, DefaultSearchDepth),
		Topic:                    defaultString(opts.Topic, DefaultTopic),
		TimeRange:                opts.TimeRange,
		Days:                     opts.Days,
		MaxResults:               defaultInt(opts.MaxResults, DefaultMaxResults),
		IncludeDomains:           opts.IncludeDomains,
		ExcludeDomains:           opts.ExcludeDomains,
		IncludeAnswer:            opts.IncludeAnswer,
		IncludeRawContent:        opts.IncludeRawContent,
		IncludeImages:            opts.IncludeImages,
		IncludeImageDescriptions: opts.IncludeImageDescriptions,
		MaxTokens:                opts.MaxTokens,
		ChunksPerSource:          opts.ChunksPerSource,
		Country:                  opts.Country,
		Timeout:                  defaultInt(opts.Timeout, 60),
	}

	var resp SearchResponse
	if err := c.doRequest(ctx, "/search", req, &resp); err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return &resp, nil
}

// Extract extracts and processes content from one or more specified URLs.
func (c *Client) Extract(ctx context.Context, urls []string, opts *ExtractOptions) (*ExtractResponse, error) {
	if len(urls) == 0 {
		return nil, &APIError{
			StatusCode: 400,
			Message:    "at least one URL is required",
		}
	}

	if opts == nil {
		opts = &ExtractOptions{}
	}

	req := &ExtractRequest{
		URLs:          urls,
		IncludeImages: opts.IncludeImages,
		ExtractDepth:  defaultString(opts.ExtractDepth, DefaultSearchDepth),
		Format:        defaultString(opts.Format, DefaultFormat),
		Timeout:       defaultInt(opts.Timeout, 60),
	}

	var resp ExtractResponse
	if err := c.doRequest(ctx, "/extract", req, &resp); err != nil {
		return nil, fmt.Errorf("extract failed: %w", err)
	}

	return &resp, nil
}

// Crawl intelligently crawls a website to discover and extract content from multiple pages.
func (c *Client) Crawl(ctx context.Context, url string, opts *CrawlOptions) (*CrawlResponse, error) {
	if url == "" {
		return nil, &APIError{
			StatusCode: 400,
			Message:    "URL is required",
		}
	}

	if opts == nil {
		opts = &CrawlOptions{}
	}

	req := &CrawlRequest{
		URL:            url,
		MaxDepth:       defaultInt(opts.MaxDepth, 1),
		MaxBreadth:     defaultInt(opts.MaxBreadth, 20),
		Limit:          defaultInt(opts.Limit, 50),
		Instructions:   opts.Instructions,
		ExtractDepth:   defaultString(opts.ExtractDepth, DefaultSearchDepth),
		SelectPaths:    opts.SelectPaths,
		SelectDomains:  opts.SelectDomains,
		ExcludePaths:   opts.ExcludePaths,
		ExcludeDomains: opts.ExcludeDomains,
		AllowExternal:  opts.AllowExternal,
		IncludeImages:  opts.IncludeImages,
		Categories:     opts.Categories,
		Format:         defaultString(opts.Format, DefaultFormat),
		Timeout:        defaultInt(opts.Timeout, 60),
	}

	var resp CrawlResponse
	if err := c.doRequest(ctx, "/crawl", req, &resp); err != nil {
		return nil, fmt.Errorf("crawl failed: %w", err)
	}

	return &resp, nil
}

// Map discovers and maps the structure of a website without extracting full content.
func (c *Client) Map(ctx context.Context, url string, opts *MapOptions) (*MapResponse, error) {
	if url == "" {
		return nil, &APIError{
			StatusCode: 400,
			Message:    "URL is required",
		}
	}

	if opts == nil {
		opts = &MapOptions{}
	}

	req := &MapRequest{
		URL:            url,
		MaxDepth:       defaultInt(opts.MaxDepth, 1),
		MaxBreadth:     defaultInt(opts.MaxBreadth, 20),
		Limit:          defaultInt(opts.Limit, 50),
		Instructions:   opts.Instructions,
		SelectPaths:    opts.SelectPaths,
		SelectDomains:  opts.SelectDomains,
		ExcludePaths:   opts.ExcludePaths,
		ExcludeDomains: opts.ExcludeDomains,
		AllowExternal:  opts.AllowExternal,
		Categories:     opts.Categories,
		Timeout:        defaultInt(opts.Timeout, 60),
	}

	var resp MapResponse
	if err := c.doRequest(ctx, "/map", req, &resp); err != nil {
		return nil, fmt.Errorf("map failed: %w", err)
	}

	return &resp, nil
}

func defaultString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func defaultInt(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}
