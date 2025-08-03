package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Searcher defines the interface for web search
type Searcher interface {
	Search(query string) (string, error)
}

// WebSearcher handles web search operations
type WebSearcher struct {
	braveAPIKey string
}

// NewWebSearcher creates a new web searcher
func NewWebSearcher(braveAPIKey string) *WebSearcher {
	return &WebSearcher{
		braveAPIKey: braveAPIKey,
	}
}

// Search performs a web search using available methods
func (s *WebSearcher) Search(query string) (string, error) {
	// Try Brave first if API key is available
	if s.braveAPIKey != "" {
		results, err := s.searchBrave(query)
		if err == nil {
			return results, nil
		}
		fmt.Printf("⚠️  Brave search failed: %v, falling back to DuckDuckGo\n", err)
	}

	// Fall back to DuckDuckGo
	return s.searchDuckDuckGo(query)
}

// searchDuckDuckGo performs a web search using DuckDuckGo
func (s *WebSearcher) searchDuckDuckGo(query string) (string, error) {
	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("DuckDuckGo search failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("DuckDuckGo returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return fmt.Sprintf("DuckDuckGo search results for '%s':\n\n%s", query, string(body)), nil
}

// braveSearchResult represents a search result from Brave Search API
type braveSearchResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// braveSearchResponse represents the response from Brave Search API
type braveSearchResponse struct {
	Web struct {
		Results []braveSearchResult `json:"results"`
	} `json:"web"`
}

// searchBrave performs a web search using Brave Search API
func (s *WebSearcher) searchBrave(query string) (string, error) {
	if s.braveAPIKey == "" {
		return "", fmt.Errorf("Brave Search API key is required")
	}

	apiURL := fmt.Sprintf("https://api.search.brave.com/res/v1/web/search?q=%s", url.QueryEscape(query))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("X-Subscription-Token", s.braveAPIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to search Brave: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Brave API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var searchResp braveSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	var results []string
	for i, result := range searchResp.Web.Results {
		if i >= 5 {
			break
		}
		results = append(results, fmt.Sprintf("Title: %s\nURL: %s\nSnippet: %s",
			result.Title, result.URL, result.Description))
	}

	if len(results) == 0 {
		return "No results found", nil
	}

	return strings.Join(results, "\n\n"), nil
}
