package tavily

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
		opts   *Options
		want   string
	}{
		{
			name:   "with api key",
			apiKey: "tvly-test-key",
			opts:   nil,
			want:   "tvly-test-key",
		},
		{
			name:   "empty api key uses env",
			apiKey: "",
			opts:   nil,
			want:   "",
		},
		{
			name:   "with custom options",
			apiKey: "tvly-test-key",
			opts: &Options{
				BaseURL: "https://custom.api.com",
				Timeout: 30 * time.Second,
			},
			want: "tvly-test-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(tt.apiKey, tt.opts)
			if client.apiKey != tt.want {
				t.Errorf("New() apiKey = %v, want %v", client.apiKey, tt.want)
			}
			if tt.opts != nil && tt.opts.BaseURL != "" {
				if client.baseURL != tt.opts.BaseURL {
					t.Errorf("New() baseURL = %v, want %v", client.baseURL, tt.opts.BaseURL)
				}
			}
		})
	}
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		checkFunc  func(*APIError) bool
	}{
		{
			name:       "unauthorized error",
			statusCode: 401,
			message:    "Invalid API key",
			checkFunc:  (*APIError).IsUnauthorized,
		},
		{
			name:       "rate limit error",
			statusCode: 429,
			message:    "Rate limit exceeded",
			checkFunc:  (*APIError).IsRateLimit,
		},
		{
			name:       "forbidden error",
			statusCode: 403,
			message:    "Access denied",
			checkFunc:  (*APIError).IsForbidden,
		},
		{
			name:       "bad request error",
			statusCode: 400,
			message:    "Invalid parameters",
			checkFunc:  (*APIError).IsBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{
				StatusCode: tt.statusCode,
				Message:    tt.message,
			}

			if err.Error() != tt.message {
				t.Errorf("APIError.Error() = %v, want %v", err.Error(), tt.message)
			}

			if !tt.checkFunc(err) {
				t.Errorf("APIError check function returned false for status %d", tt.statusCode)
			}
		})
	}
}

func TestSearchRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if !strings.Contains(r.Header.Get("Authorization"), "Bearer") {
			t.Errorf("Expected Authorization header with Bearer token")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"query": "test query",
			"response_time": 0.5,
			"images": [],
			"results": [
				{
					"title": "Test Result",
					"url": "https://example.com",
					"content": "Test content",
					"score": 0.95
				}
			]
		}`))
	}))
	defer server.Close()

	client := New("tvly-test-key", &Options{
		BaseURL: server.URL,
	})

	ctx := context.Background()
	result, err := client.Search(ctx, "test query", nil)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if result.Query != "test query" {
		t.Errorf("Search() query = %v, want %v", result.Query, "test query")
	}

	if len(result.Results) != 1 {
		t.Errorf("Search() results count = %v, want %v", len(result.Results), 1)
	}

	if result.Results[0].Title != "Test Result" {
		t.Errorf("Search() result title = %v, want %v", result.Results[0].Title, "Test Result")
	}
}

func TestSearchWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"query": "test query",
			"answer": "Test answer",
			"response_time": 0.5,
			"images": [
				{
					"url": "https://example.com/image.jpg",
					"description": "Test image"
				}
			],
			"results": []
		}`))
	}))
	defer server.Close()

	client := New("tvly-test-key", &Options{
		BaseURL: server.URL,
	})

	opts := &SearchOptions{
		SearchDepth:   string(SearchDepthAdvanced),
		Topic:         string(TopicNews),
		MaxResults:    10,
		IncludeAnswer: true,
		IncludeImages: BoolPtr(true),
	}

	ctx := context.Background()
	result, err := client.Search(ctx, "test query", opts)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if result.Answer != "Test answer" {
		t.Errorf("Search() answer = %v, want %v", result.Answer, "Test answer")
	}

	if len(result.Images) != 1 {
		t.Errorf("Search() images count = %v, want %v", len(result.Images), 1)
	}
}

func TestExtractRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"response_time": 0.5,
			"results": [
				{
					"url": "https://example.com",
					"raw_content": "Test content",
					"images": ["https://example.com/image.jpg"]
				}
			],
			"failed_results": []
		}`))
	}))
	defer server.Close()

	client := New("tvly-test-key", &Options{
		BaseURL: server.URL,
	})

	urls := []string{"https://example.com"}
	ctx := context.Background()
	result, err := client.Extract(ctx, urls, nil)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if len(result.Results) != 1 {
		t.Errorf("Extract() results count = %v, want %v", len(result.Results), 1)
	}

	if result.Results[0].URL != "https://example.com" {
		t.Errorf("Extract() result URL = %v, want %v", result.Results[0].URL, "https://example.com")
	}
}

func TestConvenienceMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"query": "test",
			"response_time": 0.5,
			"images": [],
			"results": [{"title": "Test", "url": "https://example.com", "content": "Test content", "score": 0.9}]
		}`))
	}))
	defer server.Close()

	client := New("tvly-test-key", &Options{
		BaseURL: server.URL,
	})

	ctx := context.Background()

	t.Run("SearchSimple", func(t *testing.T) {
		result, err := client.SearchSimple(ctx, "test")
		if err != nil {
			t.Fatalf("SearchSimple() error = %v", err)
		}
		if result.Query != "test" {
			t.Errorf("SearchSimple() query = %v, want %v", result.Query, "test")
		}
	})

	t.Run("SearchWithAnswer", func(t *testing.T) {
		result, err := client.SearchWithAnswer(ctx, "test")
		if err != nil {
			t.Fatalf("SearchWithAnswer() error = %v", err)
		}
		if result.Query != "test" {
			t.Errorf("SearchWithAnswer() query = %v, want %v", result.Query, "test")
		}
	})

	t.Run("SearchNews", func(t *testing.T) {
		result, err := client.SearchNews(ctx, "test", 7)
		if err != nil {
			t.Fatalf("SearchNews() error = %v", err)
		}
		if result.Query != "test" {
			t.Errorf("SearchNews() query = %v, want %v", result.Query, "test")
		}
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("BoolPtr", func(t *testing.T) {
		val := BoolPtr(true)
		if val == nil {
			t.Error("BoolPtr() returned nil")
		}
		if *val != true {
			t.Errorf("BoolPtr() = %v, want %v", *val, true)
		}
	})

	t.Run("GetVersionInfo", func(t *testing.T) {
		info := GetVersionInfo()
		if info["client_name"] != "go-tavily" {
			t.Errorf("GetVersionInfo() client_name = %v, want %v", info["client_name"], "go-tavily")
		}
		if info["client_version"] == "" {
			t.Error("GetVersionInfo() client_version is empty")
		}
	})
}

func TestErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"detail": {
				"error": "Invalid API key provided"
			}
		}`))
	}))
	defer server.Close()

	client := New("invalid-key", &Options{
		BaseURL: server.URL,
	})

	ctx := context.Background()
	_, err := client.Search(ctx, "test", nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected *APIError, got %T", err)
	}

	if !apiErr.IsUnauthorized() {
		t.Error("Expected unauthorized error")
	}

	if !strings.Contains(apiErr.Message, "Invalid API key") {
		t.Errorf("Expected error message to contain 'Invalid API key', got %v", apiErr.Message)
	}
}

func TestInputValidation(t *testing.T) {
	client := New("tvly-test-key", nil)
	ctx := context.Background()

	t.Run("Extract with empty URLs", func(t *testing.T) {
		_, err := client.Extract(ctx, []string{}, nil)
		if err == nil {
			t.Error("Expected error for empty URLs, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("Expected *APIError, got %T", err)
		}

		if !apiErr.IsBadRequest() {
			t.Error("Expected bad request error")
		}
	})

	t.Run("Crawl with empty URL", func(t *testing.T) {
		_, err := client.Crawl(ctx, "", nil)
		if err == nil {
			t.Error("Expected error for empty URL, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("Expected *APIError, got %T", err)
		}

		if !apiErr.IsBadRequest() {
			t.Error("Expected bad request error")
		}
	})

	t.Run("Map with empty URL", func(t *testing.T) {
		_, err := client.Map(ctx, "", nil)
		if err == nil {
			t.Error("Expected error for empty URL, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("Expected *APIError, got %T", err)
		}

		if !apiErr.IsBadRequest() {
			t.Error("Expected bad request error")
		}
	})
}

func BenchmarkSearch(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"query": "test", "response_time": 0.5, "images": [], "results": []}`))
	}))
	defer server.Close()

	client := New("tvly-test-key", &Options{
		BaseURL: server.URL,
	})

	ctx := context.Background()

	for b.Loop() {
		_, err := client.SearchSimple(ctx, "benchmark test")
		if err != nil {
			b.Fatal(err)
		}
	}
}
