package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GistFile represents a file in a Gist
type GistFile struct {
	Content string `json:"content"`
}

// Gist represents the Gist structure for creation/update
type Gist struct {
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
}

// GistResponse represents the response from GitHub API
type GistResponse struct {
	ID        string              `json:"id"`
	HTMLURL   string              `json:"html_url"`
	Files     map[string]GistFile `json:"files"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// Client handles GitHub Gist operations
type Client struct {
	Token string
}

// NewClient creates a new Gist client
func NewClient(token string) *Client {
	return &Client{Token: token}
}

// CreateGist creates a new Gist with the provided files
func (c *Client) CreateGist(description string, files map[string]string) (*GistResponse, error) {
	gistFiles := make(map[string]GistFile)
	for name, content := range files {
		gistFiles[name] = GistFile{Content: content}
	}

	payload := Gist{
		Description: description,
		Public:      false, // Private by default for config
		Files:       gistFiles,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/gists", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// UpdateGist updates an existing Gist
func (c *Client) UpdateGist(gistID string, files map[string]string) (*GistResponse, error) {
	gistFiles := make(map[string]GistFile)
	for name, content := range files {
		gistFiles[name] = GistFile{Content: content}
	}

	payload := Gist{
		Files: gistFiles,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://api.github.com/gists/%s", gistID), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// GetGist retrieves a Gist by ID
func (c *Client) GetGist(gistID string) (*GistResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/gists/%s", gistID), nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) (*GistResponse, error) {
	req.Header.Set("Authorization", "token "+c.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github api error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var gistResp GistResponse
	if err := json.Unmarshal(body, &gistResp); err != nil {
		return nil, fmt.Errorf("failed to parse gist response: %w", err)
	}

	return &gistResp, nil
}
