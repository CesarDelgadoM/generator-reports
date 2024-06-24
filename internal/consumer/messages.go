package consumer

import (
	"encoding/json"

	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
)

// Message queue names
type MessageQueueNames struct {
	ReportType string
	QueueName  string
}

func UnmarshalMessageQueueNames(m []byte) *MessageQueueNames {
	var msg MessageQueueNames

	if err := json.Unmarshal(m, &msg); err != nil {
		zap.Log.Error("Failed to make unmarshal to message: ", err)
		return nil
	}

	return &msg
}

// Message to consumer
type Message struct {
	UserId uint
	Format string
	Type   string
	Status int
	Data   []byte
}

func UnmarshalMessage(m []byte) *Message {
	var msg Message

	if err := json.Unmarshal(m, &msg); err != nil {
		zap.Log.Error("Failed to make unmarshal to message: ", err)
		return nil
	}

	return &msg
}
