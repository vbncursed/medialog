package bootstrap

import (
	"github.com/vbncursed/medialog/metadata_service/internal/services/metadata_service"
	library_entry_processor "github.com/vbncursed/medialog/metadata_service/internal/services/processors/library_entry_processor"
)

func InitLibraryEntryProcessor(metadataService *metadata_service.MetadataService) *library_entry_processor.LibraryEntryProcessor {
	return library_entry_processor.NewLibraryEntryProcessor(metadataService)
}

