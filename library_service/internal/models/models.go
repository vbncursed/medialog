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

// EntryStatus представляет статус записи
type EntryStatus int

const (
	EntryStatusUnspecified EntryStatus = iota
	EntryStatusPlanned
	EntryStatusInProgress
	EntryStatusDone
)

// Entry представляет запись пользователя о контенте
type Entry struct {
	EntryID    uint64
	UserID     uint64
	MediaID    uint64
	Type       MediaType
	Status     EntryStatus
	Rating     uint32
	Review     string
	Tags       []string
	StartedAt  *time.Time
	FinishedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CreateEntryInput входные данные для создания записи
type CreateEntryInput struct {
	UserID     uint64
	MediaID    uint64
	Type       MediaType
	Status     EntryStatus
	Rating     uint32
	Review     string
	Tags       []string
	StartedAt  int64
	FinishedAt int64
}

// UpdateEntryInput входные данные для обновления записи
type UpdateEntryInput struct {
	EntryID    uint64
	UserID     uint64
	Status     *EntryStatus
	Rating     *uint32
	Review     *string
	Tags       []string
	StartedAt  *int64
	FinishedAt *int64
}

// ListEntriesInput входные данные для получения списка записей
type ListEntriesInput struct {
	UserID       uint64
	Types        []MediaType
	Statuses     []EntryStatus
	Tags         []string
	MinRating    uint32
	MaxRating    uint32
	FinishedFrom int64
	FinishedTo   int64
	SortBy       string
	SortOrder    string
	Page         uint32
	PageSize     uint32
}

// ListEntriesResult результат получения списка записей
type ListEntriesResult struct {
	Entries  []*Entry
	Total    uint32
	Page     uint32
	PageSize uint32
}
