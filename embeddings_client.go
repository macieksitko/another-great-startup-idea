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
}

// NewClient creates a new API client
func NewEmbeddingsClient(baseURL string) *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        baseURL: baseURL,
    }
}

// Post represents the data structure for a post
type EmbedRequest struct {
    Texts []string `json:"texts"`
}

type EmbedResponse struct {
    Embeddings [][]float32 `json:"embeddings"`
}

// CreatePost creates a new post
func (c *Client) Embed(request EmbedRequest) ([][]float32, error) {
    jsonData, err := json.Marshal(request)
    if err != nil {
        return nil, err
    }

    resp, err := c.makeRequest("POST", "/embed", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }

    var response EmbedResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, err
    }

    return response.Embeddings, nil
}

// makeRequest is a helper function to make HTTP requests
func (c *Client) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
    req, err := http.NewRequest(method, c.baseURL+endpoint, body)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")

    return c.httpClient.Do(req)
}