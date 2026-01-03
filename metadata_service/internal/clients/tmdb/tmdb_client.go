package tmdb

import (
	"net/http"
	"time"
)

type TMDBClient struct {
	apiKey       string
	baseURL      string
	imageBaseURL string
	httpClient   *http.Client
	genreMap     map[int]string
}

func NewTMDBClient(apiKey, baseURL, imageBaseURL string, timeoutSeconds int) *TMDBClient {
	return &TMDBClient{
		apiKey:       apiKey,
		baseURL:      baseURL,
		imageBaseURL: imageBaseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
		genreMap: defaultGenreMap,
	}
}
