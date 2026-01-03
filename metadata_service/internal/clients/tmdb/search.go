package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (c *TMDBClient) SearchMedia(ctx context.Context, query string, mediaType *models.MediaType) ([]*models.Media, error) {
	if c.apiKey == "" || c.apiKey == "CHANGE_ME" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	var results []*models.Media

	if mediaType == nil || *mediaType == models.MediaTypeMovie {
		movies, err := c.searchMovies(ctx, query)
		if err == nil {
			results = append(results, movies...)
		}
	}

	if mediaType == nil || *mediaType == models.MediaTypeTV {
		tvShows, err := c.searchTV(ctx, query)
		if err == nil {
			results = append(results, tvShows...)
		}
	}

	return results, nil
}

func (c *TMDBClient) searchMovies(ctx context.Context, query string) ([]*models.Media, error) {
	reqURL := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s", c.baseURL, c.apiKey, url.QueryEscape(query))

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

	var tmdbResp TMDBMovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResp); err != nil {
		return nil, err
	}

	results := make([]*models.Media, 0, len(tmdbResp.Results))
	for _, movie := range tmdbResp.Results {
		results = append(results, c.convertMovieToMedia(&movie))
	}

	return results, nil
}

func (c *TMDBClient) searchTV(ctx context.Context, query string) ([]*models.Media, error) {
	reqURL := fmt.Sprintf("%s/search/tv?api_key=%s&query=%s", c.baseURL, c.apiKey, url.QueryEscape(query))

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

	var tmdbResp TMDBTVResponse
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResp); err != nil {
		return nil, err
	}

	results := make([]*models.Media, 0, len(tmdbResp.Results))
	for _, tv := range tmdbResp.Results {
		results = append(results, c.convertTVToMedia(&tv))
	}

	return results, nil
}

