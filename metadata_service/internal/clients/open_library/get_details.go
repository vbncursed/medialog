package open_library

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (c *OpenLibraryClient) GetMediaByExternalID(ctx context.Context, source, externalID string) (*models.Media, error) {
	if source != "isbn" && source != "openlibrary" {
		return nil, fmt.Errorf("unsupported source for Open Library: %s", source)
	}

	if source == "isbn" {
		return c.getBookByISBN(ctx, externalID)
	}

	return c.getBookByKey(ctx, externalID)
}

func (c *OpenLibraryClient) getBookByISBN(ctx context.Context, isbn string) (*models.Media, error) {
	reqURL := fmt.Sprintf("%s/isbn/%s.json", c.baseURL, isbn)

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

	var book OpenLibraryBookDetails
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		return nil, err
	}

	return c.convertBookDetailsToMedia(&book), nil
}

func (c *OpenLibraryClient) getBookByKey(ctx context.Context, key string) (*models.Media, error) {
	reqURL := fmt.Sprintf("%s%s.json", c.baseURL, key)

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

	var book OpenLibraryBookDetails
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		return nil, err
	}

	return c.convertBookDetailsToMedia(&book), nil
}

