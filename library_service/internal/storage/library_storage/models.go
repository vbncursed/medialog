package library_storage

import "errors"

const (
	entriesTable = "library_entries"
)

var (
	ErrEntryAlreadyExists = errors.New("entry already exists")
)

