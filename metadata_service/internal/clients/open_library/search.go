package open_library

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (c *OpenLibraryClient) SearchMedia(ctx context.Context, query string, mediaType *models.MediaType) ([]*models.Media, error) {
	if mediaType != nil && *mediaType != models.MediaTypeBook {
		return nil, nil
	}

	return c.searchBooks(ctx, query)
}

func (c *OpenLibraryClient) searchBooks(ctx context.Context, query string) ([]*models.Media, error) {
	reqURL := fmt.Sprintf("%s/search.json?title=%s&limit=20", c.baseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Open Library API error: %d", resp.StatusCode)
	}

	var olResp OpenLibrarySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&olResp); err != nil {
		return nil, err
	}

	results := make([]*models.Media, 0, len(olResp.Docs))
	for _, book := range olResp.Docs {
		results = append(results, c.convertBookToMedia(&book))
	}

	return results, nil
}

