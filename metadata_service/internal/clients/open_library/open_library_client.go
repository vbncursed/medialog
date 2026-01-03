package open_library

import (
	"net/http"
	"time"
)

type OpenLibraryClient struct {
	baseURL      string
	coverBaseURL string
	httpClient   *http.Client
}

func NewOpenLibraryClient(baseURL, coverBaseURL string, timeoutSeconds int) *OpenLibraryClient {
	return &OpenLibraryClient{
		baseURL:      baseURL,
		coverBaseURL: coverBaseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}
