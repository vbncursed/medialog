package library_service_api

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vbncursed/medialog/library_service/internal/models"
	pb_models "github.com/vbncursed/medialog/library_service/internal/pb/models"
)

func convertMediaType(t pb_models.MediaType) models.MediaType {
	switch t {
	case pb_models.MediaType_MEDIA_TYPE_MOVIE:
		return models.MediaTypeMovie
	case pb_models.MediaType_MEDIA_TYPE_TV:
		return models.MediaTypeTV
	case pb_models.MediaType_MEDIA_TYPE_BOOK:
		return models.MediaTypeBook
	default:
		return models.MediaTypeUnspecified
	}
}

func convertMediaTypeToProto(t models.MediaType) pb_models.MediaType {
	switch t {
	case models.MediaTypeMovie:
		return pb_models.MediaType_MEDIA_TYPE_MOVIE
	case models.MediaTypeTV:
		return pb_models.MediaType_MEDIA_TYPE_TV
	case models.MediaTypeBook:
		return pb_models.MediaType_MEDIA_TYPE_BOOK
	default:
		return pb_models.MediaType_MEDIA_TYPE_UNSPECIFIED
	}
}

func convertEntryStatus(s pb_models.EntryStatus) models.EntryStatus {
	switch s {
	case pb_models.EntryStatus_ENTRY_STATUS_PLANNED:
		return models.EntryStatusPlanned
	case pb_models.EntryStatus_ENTRY_STATUS_IN_PROGRESS:
		return models.EntryStatusInProgress
	case pb_models.EntryStatus_ENTRY_STATUS_DONE:
		return models.EntryStatusDone
	default:
		return models.EntryStatusUnspecified
	}
}

func convertEntryStatusToProto(s models.EntryStatus) pb_models.EntryStatus {
	switch s {
	case models.EntryStatusPlanned:
		return pb_models.EntryStatus_ENTRY_STATUS_PLANNED
	case models.EntryStatusInProgress:
		return pb_models.EntryStatus_ENTRY_STATUS_IN_PROGRESS
	case models.EntryStatusDone:
		return pb_models.EntryStatus_ENTRY_STATUS_DONE
	default:
		return pb_models.EntryStatus_ENTRY_STATUS_UNSPECIFIED
	}
}

func convertEntryToProto(e *models.Entry) *pb_models.Entry {
	if e == nil {
		return nil
	}

	var startedAt *timestamppb.Timestamp
	if e.StartedAt != nil {
		startedAt = timestamppb.New(*e.StartedAt)
	}

	var finishedAt *timestamppb.Timestamp
	if e.FinishedAt != nil {
		finishedAt = timestamppb.New(*e.FinishedAt)
	}

	return &pb_models.Entry{
		EntryId:    e.EntryID,
		MediaId:    e.MediaID,
		Type:       convertMediaTypeToProto(e.Type),
		Status:     convertEntryStatusToProto(e.Status),
		Rating:     e.Rating,
		Review:     e.Review,
		Tags:       e.Tags,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		CreatedAt:  timestamppb.New(e.CreatedAt),
		UpdatedAt:  timestamppb.New(e.UpdatedAt),
	}
}
