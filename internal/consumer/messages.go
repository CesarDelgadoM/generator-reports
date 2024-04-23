package consumer

// Message queue names
type MessageQueueNames struct {
	ReportType string
	QueueName  string
}

// Message to consumer
type Message struct {
	Userid uint
	Format string
	Status int
	Data   []byte
}
