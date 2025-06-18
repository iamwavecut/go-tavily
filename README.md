# 🔍 Go Tavily Client

A modern, flexible Go client for the [Tavily AI-powered search and web content extraction API](https://docs.tavily.com). Built for Go 1.24+ with modern idioms and best practices.

## 🚀 Features

| Feature             | Description                                                | Status      |
| ------------------- | ---------------------------------------------------------- | ----------- |
| **Search**          | Web search with intelligent results aggregation            | ✅ Supported |
| **Extract**         | Content extraction from specific URLs                      | ✅ Supported |
| **Crawl**           | Intelligent website crawling and content mapping           | ✅ Supported |
| **Map**             | Website structure discovery and mapping                    | ✅ Supported |
| **Context Support** | Full context.Context integration for cancellation/timeouts | ✅ Modern    |
| **Error Handling**  | Typed errors with semantic checking methods                | ✅ Robust    |
| **Testing**         | Comprehensive test suite                                   | ✅ Complete  |
| **Performance**     | Built for Go 1.24+                                         | ⚡ Optimized |

## 📦 Installation

```bash
go get github.com/iamwavecut/go-tavily
```

## 🔧 Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/iamwavecut/go-tavily"
)

func main() {
    // Create client (uses TAVILY_API_KEY env var if empty)
    client := tavily.New("tvly-your-api-key", nil)
    
    ctx := context.Background()
    
    // Simple search
    result, err := client.SearchSimple(ctx, "Go programming language")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d results in %.2fs\n", 
        len(result.Results), result.ResponseTime)
}
```

## 📖 Usage Examples

### 🔍 Advanced Search

```go
opts := &tavily.SearchOptions{
    SearchDepth:    string(tavily.SearchDepthAdvanced),
    Topic:          string(tavily.TopicNews),
    MaxResults:     10,
    IncludeAnswer:  true,
    IncludeImages:  tavily.BoolPtr(true),
    TimeRange:      string(tavily.TimeRangeWeek),
    IncludeDomains: []string{"github.com", "golang.org"},
    Country:        "US",
}

result, err := client.Search(ctx, "Go 1.24 release", opts)
```

### 🌐 Content Extraction

```go
urls := []string{
    "https://golang.org/doc/",
    "https://pkg.go.dev/",
}

opts := &tavily.ExtractOptions{
    Format:        string(tavily.FormatMarkdown),
    ExtractDepth:  string(tavily.SearchDepthAdvanced),
    IncludeImages: tavily.BoolPtr(true),
}

result, err := client.Extract(ctx, urls, opts)
```

### 🕷️ Website Crawling

```go
opts := &tavily.CrawlOptions{
    MaxDepth:      2,
    MaxBreadth:    10,
    Limit:         20,
    SelectPaths:   []string{"/docs/*", "/api/*"},
    Categories:    []tavily.CrawlCategory{
        tavily.CategoryDocumentation,
        tavily.CategoryDeveloper,
    },
    Format:        string(tavily.FormatMarkdown),
    AllowExternal: tavily.BoolPtr(false),
}

result, err := client.Crawl(ctx, "https://docs.tavily.com", opts)
```

### 🗺️ Website Mapping

```go
opts := &tavily.MapOptions{
    MaxDepth:    3,
    Limit:       100,
    Categories:  []tavily.CrawlCategory{
        tavily.CategoryDocumentation,
        tavily.CategoryBlog,
    },
    SelectPaths: []string{"/docs/*", "/blog/*"},
}

result, err := client.Map(ctx, "https://docs.tavily.com", opts)
```

## 🎯 Convenience Methods

| Method                 | Purpose                              | Example          |
| ---------------------- | ------------------------------------ | ---------------- |
| `SearchSimple()`       | Basic search with minimal config     | Quick searches   |
| `SearchWithAnswer()`   | Search with AI-generated answer      | Q&A applications |
| `SearchNews()`         | News-focused search with time filter | Recent updates   |
| `ExtractSimple()`      | Single URL extraction                | Content analysis |
| `ExtractWithImages()`  | Multi-URL extraction with images     | Rich content     |
| `CrawlDocumentation()` | Documentation-focused crawling       | API docs, guides |
| `MapSite()`            | Quick website structure mapping      | Site analysis    |
| `GetSearchContext()`   | RAG-formatted search results         | AI applications  |

## 🛠️ Configuration

### Client Options

```go
opts := &tavily.Options{
    BaseURL:    "https://api.tavily.com",  // Custom API endpoint
    HTTPClient: customHTTPClient,          // Custom HTTP client
    Timeout:    45 * time.Second,         // Request timeout
}

client := tavily.New("your-api-key", opts)
```

### Custom HTTP Client

```go
customClient := &http.Client{
    Timeout: 45 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:       10,
        IdleConnTimeout:    30 * time.Second,
        DisableCompression: true,
    },
}

client := tavily.New("your-api-key", &tavily.Options{
    HTTPClient: customClient,
})
```

## 🚨 Error Handling

The client provides semantic error checking methods:

```go
result, err := client.Search(ctx, "query", nil)
if err != nil {
    if apiErr, ok := err.(*tavily.APIError); ok {
        switch {
        case apiErr.IsUnauthorized():
            fmt.Println("Invalid API key")
        case apiErr.IsRateLimit():
            fmt.Println("Rate limit exceeded")
        case apiErr.IsForbidden():
            fmt.Println("Access forbidden")
        case apiErr.IsBadRequest():
            fmt.Println("Invalid parameters")
        default:
            fmt.Printf("API error: %s\n", apiErr.Message)
        }
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## 🧪 Testing

The client includes comprehensive tests:

```bash
# Run all tests
go test -v ./...

# Run benchmarks with new testing.B.Loop
go test -bench=. -count=3

# Test with coverage
go test -cover ./...
```

## 🏃‍♂️ Demo Application

```bash
# Set your API key
export TAVILY_API_KEY="tvly-your-api-key"

# Run the demo
go run cmd/demo/main.go
```

## 🔧 Environment Variables

| Variable             | Description         | Required   |
| -------------------- | ------------------- | ---------- |
| `TAVILY_API_KEY`     | Your Tavily API key | ✅ Yes      |
| `TAVILY_HTTP_PROXY`  | HTTP proxy URL      | ❌ Optional |
| `TAVILY_HTTPS_PROXY` | HTTPS proxy URL     | ❌ Optional |

## 📋 API Coverage

| Endpoint   | Method      | Status     | Features                                  |
| ---------- | ----------- | ---------- | ----------------------------------------- |
| `/search`  | `Search()`  | ✅ Complete | All parameters, answer generation, images |
| `/extract` | `Extract()` | ✅ Complete | Multi-URL, formats, depth control         |
| `/crawl`   | `Crawl()`   | ✅ Complete | Path filtering, categories, depth limits  |
| `/map`     | `Map()`     | ✅ Complete | Structure discovery, URL filtering        |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `go test -v ./...`
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Links

- [Tavily API Documentation](https://docs.tavily.com)
