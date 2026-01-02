package library_entry_event_producer

type LibraryEntryEventProducer struct {
	kafkaBrokers []string
	topicName    string
}

func NewLibraryEntryEventProducer(kafkaBrokers []string, topicName string) *LibraryEntryEventProducer {
	return &LibraryEntryEventProducer{
		kafkaBrokers: kafkaBrokers,
		topicName:    topicName,
	}
}
