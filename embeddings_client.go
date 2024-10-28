package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles API communications
type Client struct {
    httpClient *http.Client
    baseURL    string
    apiKey     string
}

// NewClient creates a new API client
func NewClient(baseURL, apiKey string) *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        baseURL: baseURL,
        apiKey:  apiKey,
    }
}

// Post represents the data structure for a post
type Post struct {
    UserID int    `json:"userId"`
    Title  string `json:"title"`
    Body   string `json:"body"`
}

// GetPosts fetches posts from the API
func (c *Client) GetPosts() ([]Post, error) {
    resp, err := c.makeRequest("GET", "/posts", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var posts []Post
    if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
        return nil, err
    }
    return posts, nil
}

// CreatePost creates a new post
func (c *Client) CreatePost(post Post) error {
    jsonData, err := json.Marshal(post)
    if err != nil {
        return err
    }

    resp, err := c.makeRequest("POST", "/posts", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }

    return nil
}

// makeRequest is a helper function to make HTTP requests
func (c *Client) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
    req, err := http.NewRequest(method, c.baseURL+endpoint, body)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)

    return c.httpClient.Do(req)
}