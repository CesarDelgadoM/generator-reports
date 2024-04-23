package utils

import (
	"encoding/json"

	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/models"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
)

func UnmarshalMessageQueueNames(m []byte) *consumer.MessageQueueNames {
	var msg consumer.MessageQueueNames

	if err := json.Unmarshal(m, &msg); err != nil {
		zap.Log.Error("Failed to make unmarshal to message: ", err)
		return nil
	}

	return &msg
}

func UnmarshalMessage(m []byte) *consumer.Message {
	var msg consumer.Message

	if err := json.Unmarshal(m, &msg); err != nil {
		zap.Log.Error("Failed to make unmarshal to message: ", err)
		return nil
	}

	return &msg
}

func UnmarshalRestaurant(msg []byte) *models.Restaurant {
	var restaurant models.Restaurant

	if err := json.Unmarshal(msg, &restaurant); err != nil {
		zap.Log.Error("Failed to make unmarshal to restaurant: ", err)
		return nil
	}

	return &restaurant
}

func UnmarshalBranches(msg []byte) *[]models.Branch {
	var branches []models.Branch

	if err := json.Unmarshal(msg, &branches); err != nil {
		zap.Log.Error("Failed to make unmarshal to branches: ", err)
		return nil
	}

	return &branches
}
