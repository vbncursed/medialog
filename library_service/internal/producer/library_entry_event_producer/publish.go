package library_entry_event_producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/vbncursed/medialog/library_service/internal/models"
	metadata_models "github.com/vbncursed/medialog/library_service/internal/pb/metadata_models"
	"github.com/vbncursed/medialog/shared/events"
	"google.golang.org/grpc/metadata"
)

// PublishEntryChanged публикует событие об изменении записи
func (p *LibraryEntryEventProducer) PublishEntryChanged(ctx context.Context, entry *models.Entry) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(p.kafkaBrokers...),
		Topic:    p.topicName,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	event := &events.LibraryEntryEvent{
		EntryID:   entry.EntryID,
		UserID:    entry.UserID,
		MediaID:   entry.MediaID,
		Type:      int(entry.Type),
		Status:    int(entry.Status),
		Rating:    entry.Rating,
		UpdatedAt: entry.UpdatedAt.Unix(),
	}

	if p.metadataClient != nil && entry.MediaID > 0 {
		// Создаем новый контекст с JWT токеном из исходного контекста
		grpcCtx := p.addAuthToContext(ctx)

		mediaResp, err := p.metadataClient.GetMedia(grpcCtx, &metadata_models.GetMediaRequest{
			MediaId: entry.MediaID,
		})
		if err != nil {
			slog.Debug("failed to get media from metadata_service", "error", err, "media_id", entry.MediaID)
		} else if mediaResp != nil && mediaResp.Media != nil {
			if len(mediaResp.Media.ExternalIds) > 0 {
				extID := mediaResp.Media.ExternalIds[0]
				event.ExternalID = &events.ExternalIDEvent{
					Source:     extID.Source,
					ExternalID: extID.ExternalId,
				}
				slog.Debug("external_id added to event", "source", extID.Source, "external_id", extID.ExternalId, "entry_id", entry.EntryID)
			} else {
				slog.Debug("media found but no external_ids", "media_id", entry.MediaID)
			}
		}
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal entry event", "err", err)
		return err
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", entry.EntryID)),
		Value: eventJSON,
	}

	err = writer.WriteMessages(ctx, msg)
	if err != nil {
		slog.Error("failed to publish entry event", "err", err, "entry_id", entry.EntryID)
		return err
	}

	slog.Info("entry event published", "entry_id", entry.EntryID, "topic", p.topicName, "has_external_id", event.ExternalID != nil)
	return nil
}

// addAuthToContext извлекает JWT токен из входящего контекста и добавляет его в gRPC metadata
func (p *LibraryEntryEventProducer) addAuthToContext(ctx context.Context) context.Context {
	// Извлекаем токен из входящего gRPC контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		// Пробуем извлечь из обычного контекста (если это HTTP запрос через gateway)
		return ctx
	}

	// Ищем токен в заголовке Authorization
	var tokenString string
	if authHeaders := md.Get("authorization"); len(authHeaders) > 0 {
		authHeader := authHeaders[0]
		// Поддерживаем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString = parts[1]
		} else {
			tokenString = authHeader
		}
	}

	if tokenString == "" {
		return ctx
	}

	// Создаем новый gRPC metadata с токеном для исходящего вызова
	outgoingMD := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", tokenString),
	})

	return metadata.NewOutgoingContext(ctx, outgoingMD)
}
