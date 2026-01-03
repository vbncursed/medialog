package models

import "time"

const (
	RoleGuest = "guest"
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// MediaType представляет тип контента
type MediaType int

const (
	MediaTypeUnspecified MediaType = iota
	MediaTypeMovie
	MediaTypeTV
	MediaTypeBook
)

// ExternalID представляет внешний идентификатор контента
type ExternalID struct {
	Source     string // tmdb, imdb, isbn и т.д.
	ExternalID string
}

// Media представляет метаданные контента
type Media struct {
	MediaID     uint64
	Type        MediaType
	Title       string
	Year        *uint32
	Genres      []string
	PosterURL   *string // для фильмов/сериалов
	CoverURL    *string // для книг
	ExternalIDs []ExternalID
	UpdatedAt   time.Time
}

// SearchMediaInput входные данные для поиска контента
type SearchMediaInput struct {
	Query      string
	Type       *MediaType
	ExternalID *ExternalID
	Page       uint32
	PageSize   uint32
}

// SearchMediaResult результат поиска контента
type SearchMediaResult struct {
	Results  []*Media
	Total    uint32
	Page     uint32
	PageSize uint32
}

// CreateMediaInput входные данные для создания записи media
type CreateMediaInput struct {
	Type        MediaType
	Title       string
	Year        *uint32
	Genres      []string
	PosterURL   *string
	CoverURL    *string
	ExternalIDs []ExternalID
}
