package tavily

// APIError represents an error response from the Tavily API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

// IsRateLimit returns true if the error is due to rate limiting.
func (e *APIError) IsRateLimit() bool {
	return e.StatusCode == 429
}

// IsForbidden returns true if the error is due to access denied or usage limit exceeded.
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == 403 || e.StatusCode == 432 || e.StatusCode == 433
}

// IsUnauthorized returns true if the error is due to invalid API key.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsBadRequest returns true if the error is due to invalid parameters.
func (e *APIError) IsBadRequest() bool {
	return e.StatusCode == 400
}

// SearchDepth represents the depth level for search operations.
type SearchDepth string

const (
	SearchDepthBasic    SearchDepth = "basic"
	SearchDepthAdvanced SearchDepth = "advanced"
)

// Topic represents the topic category for search operations.
type Topic string

const (
	TopicGeneral Topic = "general"
	TopicNews    Topic = "news"
	TopicFinance Topic = "finance"
)

// TimeRange represents the time range for search results.
type TimeRange string

const (
	TimeRangeDay   TimeRange = "day"
	TimeRangeWeek  TimeRange = "week"
	TimeRangeMonth TimeRange = "month"
	TimeRangeYear  TimeRange = "year"
	TimeRangeD     TimeRange = "d"
	TimeRangeW     TimeRange = "w"
	TimeRangeM     TimeRange = "m"
	TimeRangeY     TimeRange = "y"
)

// Format represents the output format for content.
type Format string

const (
	FormatText     Format = "text"
	FormatMarkdown Format = "markdown"
)

// CrawlCategory represents content categories for filtering crawl results.
type CrawlCategory string

const (
	CategoryDocumentation  CrawlCategory = "Documentation"
	CategoryBlog           CrawlCategory = "Blog"
	CategoryBlogs          CrawlCategory = "Blogs"
	CategoryCommunity      CrawlCategory = "Community"
	CategoryAbout          CrawlCategory = "About"
	CategoryContact        CrawlCategory = "Contact"
	CategoryPrivacy        CrawlCategory = "Privacy"
	CategoryTerms          CrawlCategory = "Terms"
	CategoryStatus         CrawlCategory = "Status"
	CategoryPricing        CrawlCategory = "Pricing"
	CategoryEnterprise     CrawlCategory = "Enterprise"
	CategoryCareers        CrawlCategory = "Careers"
	CategoryECommerce      CrawlCategory = "E-Commerce"
	CategoryAuthentication CrawlCategory = "Authentication"
	CategoryDeveloper      CrawlCategory = "Developer"
	CategoryDevelopers     CrawlCategory = "Developers"
	CategorySolutions      CrawlCategory = "Solutions"
	CategoryPartners       CrawlCategory = "Partners"
	CategoryDownloads      CrawlCategory = "Downloads"
	CategoryMedia          CrawlCategory = "Media"
	CategoryEvents         CrawlCategory = "Events"
	CategoryPeople         CrawlCategory = "People"
)

// SearchOptions contains optional parameters for search requests.
type SearchOptions struct {
	SearchDepth              string
	Topic                    string
	TimeRange                string
	Days                     int
	MaxResults               int
	IncludeDomains           []string
	ExcludeDomains           []string
	IncludeAnswer            any
	IncludeRawContent        any
	IncludeImages            *bool
	IncludeImageDescriptions *bool
	MaxTokens                int
	ChunksPerSource          int
	Country                  string
	Timeout                  int
}

// ExtractOptions contains optional parameters for extract requests.
type ExtractOptions struct {
	IncludeImages *bool
	ExtractDepth  string
	Format        string
	Timeout       int
}

// CrawlOptions contains optional parameters for crawl requests.
type CrawlOptions struct {
	MaxDepth       int
	MaxBreadth     int
	Limit          int
	Instructions   string
	ExtractDepth   string
	SelectPaths    []string
	SelectDomains  []string
	ExcludePaths   []string
	ExcludeDomains []string
	AllowExternal  *bool
	IncludeImages  *bool
	Categories     []CrawlCategory
	Format         string
	Timeout        int
}

// MapOptions contains optional parameters for map requests.
type MapOptions struct {
	MaxDepth       int
	MaxBreadth     int
	Limit          int
	Instructions   string
	SelectPaths    []string
	SelectDomains  []string
	ExcludePaths   []string
	ExcludeDomains []string
	AllowExternal  *bool
	Categories     []CrawlCategory
	Timeout        int
}

// SearchRequest represents the request payload for search operations.
type SearchRequest struct {
	Query                    string   `json:"query"`
	SearchDepth              string   `json:"search_depth,omitempty"`
	Topic                    string   `json:"topic,omitempty"`
	TimeRange                string   `json:"time_range,omitempty"`
	Days                     int      `json:"days,omitempty"`
	MaxResults               int      `json:"max_results,omitempty"`
	IncludeDomains           []string `json:"include_domains,omitempty"`
	ExcludeDomains           []string `json:"exclude_domains,omitempty"`
	IncludeAnswer            any      `json:"include_answer,omitempty"`
	IncludeRawContent        any      `json:"include_raw_content,omitempty"`
	IncludeImages            *bool    `json:"include_images,omitempty"`
	IncludeImageDescriptions *bool    `json:"include_image_descriptions,omitempty"`
	MaxTokens                int      `json:"max_tokens,omitempty"`
	ChunksPerSource          int      `json:"chunks_per_source,omitempty"`
	Country                  string   `json:"country,omitempty"`
	Timeout                  int      `json:"timeout,omitempty"`
}

// ExtractRequest represents the request payload for extract operations.
type ExtractRequest struct {
	URLs          []string `json:"urls"`
	IncludeImages *bool    `json:"include_images,omitempty"`
	ExtractDepth  string   `json:"extract_depth,omitempty"`
	Format        string   `json:"format,omitempty"`
	Timeout       int      `json:"timeout,omitempty"`
}

// CrawlRequest represents the request payload for crawl operations.
type CrawlRequest struct {
	URL            string          `json:"url"`
	MaxDepth       int             `json:"max_depth,omitempty"`
	MaxBreadth     int             `json:"max_breadth,omitempty"`
	Limit          int             `json:"limit,omitempty"`
	Instructions   string          `json:"instructions,omitempty"`
	ExtractDepth   string          `json:"extract_depth,omitempty"`
	SelectPaths    []string        `json:"select_paths,omitempty"`
	SelectDomains  []string        `json:"select_domains,omitempty"`
	ExcludePaths   []string        `json:"exclude_paths,omitempty"`
	ExcludeDomains []string        `json:"exclude_domains,omitempty"`
	AllowExternal  *bool           `json:"allow_external,omitempty"`
	IncludeImages  *bool           `json:"include_images,omitempty"`
	Categories     []CrawlCategory `json:"categories,omitempty"`
	Format         string          `json:"format,omitempty"`
	Timeout        int             `json:"timeout,omitempty"`
}

// MapRequest represents the request payload for map operations.
type MapRequest struct {
	URL            string          `json:"url"`
	MaxDepth       int             `json:"max_depth,omitempty"`
	MaxBreadth     int             `json:"max_breadth,omitempty"`
	Limit          int             `json:"limit,omitempty"`
	Instructions   string          `json:"instructions,omitempty"`
	SelectPaths    []string        `json:"select_paths,omitempty"`
	SelectDomains  []string        `json:"select_domains,omitempty"`
	ExcludePaths   []string        `json:"exclude_paths,omitempty"`
	ExcludeDomains []string        `json:"exclude_domains,omitempty"`
	AllowExternal  *bool           `json:"allow_external,omitempty"`
	Categories     []CrawlCategory `json:"categories,omitempty"`
	Timeout        int             `json:"timeout,omitempty"`
}

// SearchImage represents an image found in search results.
type SearchImage struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// SearchResult represents a single search result.
type SearchResult struct {
	Title         string  `json:"title"`
	URL           string  `json:"url"`
	Content       string  `json:"content"`
	RawContent    string  `json:"raw_content,omitempty"`
	Score         float64 `json:"score"`
	PublishedDate string  `json:"published_date,omitempty"`
}

// SearchResponse represents the response from search operations.
type SearchResponse struct {
	Query        string         `json:"query"`
	Answer       string         `json:"answer,omitempty"`
	ResponseTime float64        `json:"response_time"`
	Images       []SearchImage  `json:"images"`
	Results      []SearchResult `json:"results"`
}

// ExtractResult represents a successful content extraction.
type ExtractResult struct {
	URL        string   `json:"url"`
	RawContent string   `json:"raw_content"`
	Images     []string `json:"images,omitempty"`
}

// ExtractFailedResult represents a failed content extraction.
type ExtractFailedResult struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

// ExtractResponse represents the response from extract operations.
type ExtractResponse struct {
	ResponseTime  float64               `json:"response_time"`
	Results       []ExtractResult       `json:"results"`
	FailedResults []ExtractFailedResult `json:"failed_results"`
}

// CrawlResult represents a crawled page with content.
type CrawlResult struct {
	URL        string   `json:"url"`
	RawContent string   `json:"raw_content"`
	Images     []string `json:"images,omitempty"`
}

// CrawlResponse represents the response from crawl operations.
type CrawlResponse struct {
	ResponseTime float64       `json:"response_time"`
	BaseURL      string        `json:"base_url"`
	Results      []CrawlResult `json:"results"`
}

// MapResponse represents the response from map operations.
type MapResponse struct {
	ResponseTime float64  `json:"response_time"`
	BaseURL      string   `json:"base_url"`
	Results      []string `json:"results"`
}
