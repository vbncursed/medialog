package events

// LibraryEntryEvent представляет событие об изменении записи в библиотеке
type LibraryEntryEvent struct {
	EntryID    uint64           `json:"entry_id"`
	UserID     uint64           `json:"user_id"`
	MediaID    uint64           `json:"media_id"`
	Type       int              `json:"type"`
	Status     int              `json:"status"`
	Rating     uint32           `json:"rating"`
	UpdatedAt  int64            `json:"updated_at"`
	ExternalID *ExternalIDEvent `json:"external_id,omitempty"`
}

// ExternalIDEvent представляет внешний идентификатор в событии
type ExternalIDEvent struct {
	Source     string `json:"source"`
	ExternalID string `json:"external_id"`
}
