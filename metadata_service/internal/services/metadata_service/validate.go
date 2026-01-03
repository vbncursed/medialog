package metadata_service

import (
	"errors"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (s *MetadataService) validateSearchInput(input models.SearchMediaInput) error {
	if input.Query == "" && input.ExternalID == nil {
		return errors.New("query or external_id is required")
	}

	if input.Page == 0 {
		return errors.New("page must be greater than 0")
	}

	if input.PageSize == 0 {
		return errors.New("page_size must be greater than 0")
	}

	if input.Type != nil {
		if *input.Type < models.MediaTypeMovie || *input.Type > models.MediaTypeBook {
			return ErrInvalidMediaType
		}
	}

	return nil
}

func (s *MetadataService) validateCreateInput(input models.CreateMediaInput) error {
	if input.Title == "" {
		return errors.New("title is required")
	}

	if input.Type < models.MediaTypeMovie || input.Type > models.MediaTypeBook {
		return ErrInvalidMediaType
	}

	if input.Year != nil && *input.Year > 3000 {
		return errors.New("year is invalid")
	}

	return nil
}

