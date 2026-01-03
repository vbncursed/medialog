package metadata_service_api

import (
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ErrCodeInvalidInput = "INVALID_INPUT"
	ErrCodeMissingField = "MISSING_FIELD"

	ErrCodeMediaNotFound = "MEDIA_NOT_FOUND"

	ErrCodeInvalidMediaID   = "INVALID_MEDIA_ID"
	ErrCodeInvalidMediaType = "INVALID_MEDIA_TYPE"
	ErrCodeInvalidSource    = "INVALID_SOURCE"

	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeInternal     = "INTERNAL_ERROR"
)

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Field   string            `json:"field,omitempty"`
	Meta    map[string]string `json:"meta,omitempty"`
}

func newError(grpcCode codes.Code, errCode, message string) error {
	detail := ErrorDetail{
		Code:    errCode,
		Message: message,
	}
	jsonBytes, err := json.Marshal(detail)
	if err != nil {
		return status.Error(grpcCode, message)
	}
	return status.Error(grpcCode, string(jsonBytes))
}

func newFieldError(grpcCode codes.Code, errCode, field, message string) error {
	detail := ErrorDetail{
		Code:    errCode,
		Message: message,
		Field:   field,
	}
	jsonBytes, err := json.Marshal(detail)
	if err != nil {
		return status.Error(grpcCode, message)
	}
	return status.Error(grpcCode, string(jsonBytes))
}

