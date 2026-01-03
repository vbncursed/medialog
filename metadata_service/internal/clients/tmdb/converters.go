package tmdb

import (
	"fmt"
	"time"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (c *TMDBClient) convertMovieToMedia(movie interface{}) *models.Media {
	var title string
	var releaseDate string
	var posterPath *string
	var genres []string
	var externalID string

	switch m := movie.(type) {
	case *TMDBMovie:
		title = m.Title
		releaseDate = m.ReleaseDate
		posterPath = m.PosterPath
		externalID = fmt.Sprintf("%d", m.ID)
		for _, genreID := range m.GenreIDs {
			if genreName, ok := c.genreMap[genreID]; ok {
				genres = append(genres, genreName)
			}
		}
	case *TMDBMovieDetails:
		title = m.Title
		releaseDate = m.ReleaseDate
		posterPath = m.PosterPath
		externalID = fmt.Sprintf("%d", m.ID)
		for _, g := range m.Genres {
			genres = append(genres, g.Name)
		}
	}

	var year *uint32
	if len(releaseDate) >= 4 {
		var y uint32
		if _, err := fmt.Sscanf(releaseDate[:4], "%d", &y); err == nil {
			year = &y
		}
	}

	var posterURL *string
	if posterPath != nil && *posterPath != "" {
		url := c.imageBaseURL + *posterPath
		posterURL = &url
	}

	if genres == nil {
		genres = []string{}
	}
	return &models.Media{
		Type:      models.MediaTypeMovie,
		Title:     title,
		Year:      year,
		Genres:    genres,
		PosterURL: posterURL,
		ExternalIDs: []models.ExternalID{
			{Source: "tmdb", ExternalID: externalID},
		},
		UpdatedAt: time.Now(),
	}
}

func (c *TMDBClient) convertTVToMedia(tv interface{}) *models.Media {
	var title string
	var firstAirDate string
	var posterPath *string
	var genres []string
	var externalID string

	switch t := tv.(type) {
	case *TMDBTV:
		title = t.Name
		firstAirDate = t.FirstAirDate
		posterPath = t.PosterPath
		externalID = fmt.Sprintf("%d", t.ID)
		for _, genreID := range t.GenreIDs {
			if genreName, ok := c.genreMap[genreID]; ok {
				genres = append(genres, genreName)
			}
		}
	case *TMDBTVDetails:
		title = t.Name
		firstAirDate = t.FirstAirDate
		posterPath = t.PosterPath
		externalID = fmt.Sprintf("%d", t.ID)
		for _, g := range t.Genres {
			genres = append(genres, g.Name)
		}
	}

	var year *uint32
	if len(firstAirDate) >= 4 {
		var y uint32
		if _, err := fmt.Sscanf(firstAirDate[:4], "%d", &y); err == nil {
			year = &y
		}
	}

	var posterURL *string
	if posterPath != nil && *posterPath != "" {
		url := c.imageBaseURL + *posterPath
		posterURL = &url
	}

	if genres == nil {
		genres = []string{}
	}
	return &models.Media{
		Type:      models.MediaTypeTV,
		Title:     title,
		Year:      year,
		Genres:    genres,
		PosterURL: posterURL,
		ExternalIDs: []models.ExternalID{
			{Source: "tmdb", ExternalID: externalID},
		},
		UpdatedAt: time.Now(),
	}
}
