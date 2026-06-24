package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HttpClient interface {
	SendJSON(path string, data interface{}, out interface{}) error
	FetchJSON(path string, out interface{}) error
	FetchJSONWithAuth(path string, bearerToken string, out interface{}) error
}

type httpClientImpl struct {
	Client  *http.Client
	BaseURL string
}

func NewHttpClient(baseURL string) HttpClient {
	return &httpClientImpl{
		Client:  &http.Client{},
		BaseURL: baseURL,
	}
}

func NewHttpClientWithCustomClient(baseURL string, client *http.Client) HttpClient {
	return &httpClientImpl{
		Client:  client,
		BaseURL: baseURL,
	}
}

func (c *httpClientImpl) SendJSON(path string, data interface{}, out interface{}) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	fmt.Printf("Sending JSON to %s: %+v", url, data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// TODO: Retry?
	// TODO: Make this not only POST...
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to decode response JSON: %w", err)
		}
	}

	return nil
}

func (c *httpClientImpl) FetchJSON(path string, out interface{}) error {
	return c.fetchJSON(path, "", out)
}

func (c *httpClientImpl) FetchJSONWithAuth(path string, bearerToken string, out interface{}) error {
	return c.fetchJSON(path, bearerToken, out)
}

func (c *httpClientImpl) fetchJSON(path string, bearerToken string, out interface{}) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	log.Printf("Fetching JSON from %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}
