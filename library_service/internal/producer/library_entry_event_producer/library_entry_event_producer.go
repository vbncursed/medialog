package library_entry_event_producer

import (
	"log/slog"

	metadata_api "github.com/vbncursed/medialog/library_service/internal/pb/metadata_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LibraryEntryEventProducer struct {
	kafkaBrokers        []string
	topicName           string
	metadataServiceAddr string
	grpcConn            *grpc.ClientConn
	metadataClient      metadata_api.MetadataServiceClient
}

func NewLibraryEntryEventProducer(kafkaBrokers []string, topicName string, metadataServiceAddr string) *LibraryEntryEventProducer {
	producer := &LibraryEntryEventProducer{
		kafkaBrokers:        kafkaBrokers,
		topicName:           topicName,
		metadataServiceAddr: metadataServiceAddr,
	}

	if metadataServiceAddr != "" {
		conn, err := grpc.NewClient(metadataServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			slog.Warn("failed to create gRPC connection to metadata_service", "error", err, "addr", metadataServiceAddr)
		} else {
			producer.grpcConn = conn
			producer.metadataClient = metadata_api.NewMetadataServiceClient(conn)
		}
	}

	return producer
}

func (p *LibraryEntryEventProducer) Close() error {
	if p.grpcConn != nil {
		return p.grpcConn.Close()
	}
	return nil
}
