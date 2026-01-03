package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (c *TMDBClient) GetMediaByExternalID(ctx context.Context, source, externalID string) (*models.Media, error) {
	if source != "tmdb" {
		return nil, fmt.Errorf("unsupported source for TMDB: %s", source)
	}

	movie, err := c.getMovieDetails(ctx, externalID)
	if err == nil && movie != nil {
		return c.convertMovieToMedia(movie), nil
	}

	tv, err := c.getTVDetails(ctx, externalID)
	if err == nil && tv != nil {
		return c.convertTVToMedia(tv), nil
	}

	return nil, fmt.Errorf("media not found")
}

func (c *TMDBClient) getMovieDetails(ctx context.Context, movieID string) (*TMDBMovieDetails, error) {
	reqURL := fmt.Sprintf("%s/movie/%s?api_key=%s", c.baseURL, movieID, c.apiKey)

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
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	var movie TMDBMovieDetails
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

func (c *TMDBClient) getTVDetails(ctx context.Context, tvID string) (*TMDBTVDetails, error) {
	reqURL := fmt.Sprintf("%s/tv/%s?api_key=%s", c.baseURL, tvID, c.apiKey)

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
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	var tv TMDBTVDetails
	if err := json.NewDecoder(resp.Body).Decode(&tv); err != nil {
		return nil, err
	}

	return &tv, nil
}

