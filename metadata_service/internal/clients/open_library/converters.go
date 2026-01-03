package open_library

import (
	"fmt"
	"time"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (c *OpenLibraryClient) convertBookToMedia(book *OpenLibraryBook) *models.Media {
	var year *uint32
	if book.FirstPublishYear != nil {
		y := uint32(*book.FirstPublishYear)
		year = &y
	}

	var coverURL *string
	if book.CoverI != nil {
		url := fmt.Sprintf("%s/id/%d-L.jpg", c.coverBaseURL, *book.CoverI)
		coverURL = &url
	}

	externalIDs := []models.ExternalID{
		{Source: "openlibrary", ExternalID: book.Key},
	}
	if len(book.ISBN) > 0 {
		externalIDs = append(externalIDs, models.ExternalID{
			Source:     "isbn",
			ExternalID: book.ISBN[0],
		})
	}

	genres := make([]string, 0)
	if len(book.Subject) > 0 {
		genres = book.Subject[:min(5, len(book.Subject))]
	}
	if genres == nil {
		genres = []string{}
	}

	return &models.Media{
		Type:        models.MediaTypeBook,
		Title:       book.Title,
		Year:        year,
		Genres:      genres,
		CoverURL:    coverURL,
		ExternalIDs: externalIDs,
		UpdatedAt:   time.Now(),
	}
}

func (c *OpenLibraryClient) convertBookDetailsToMedia(book *OpenLibraryBookDetails) *models.Media {
	var year *uint32
	if book.FirstPublishYear != nil {
		y := uint32(*book.FirstPublishYear)
		year = &y
	}

	var coverURL *string
	if len(book.Covers) > 0 {
		url := fmt.Sprintf("%s/id/%d-L.jpg", c.coverBaseURL, book.Covers[0])
		coverURL = &url
	}

	externalIDs := []models.ExternalID{
		{Source: "openlibrary", ExternalID: book.Key},
	}
	if len(book.ISBN10) > 0 {
		externalIDs = append(externalIDs, models.ExternalID{
			Source:     "isbn",
			ExternalID: book.ISBN10[0],
		})
	} else if len(book.ISBN13) > 0 {
		externalIDs = append(externalIDs, models.ExternalID{
			Source:     "isbn",
			ExternalID: book.ISBN13[0],
		})
	}

	genres := make([]string, 0)
	if len(book.Subjects) > 0 {
		genres = book.Subjects[:min(5, len(book.Subjects))]
	}
	if genres == nil {
		genres = []string{}
	}

	return &models.Media{
		Type:        models.MediaTypeBook,
		Title:       book.Title,
		Year:        year,
		Genres:      genres,
		CoverURL:    coverURL,
		ExternalIDs: externalIDs,
		UpdatedAt:   time.Now(),
	}
}

