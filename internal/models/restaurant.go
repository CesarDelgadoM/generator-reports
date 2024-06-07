package models

import (
	"encoding/json"

	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
)

// Model restaurant data
type Restaurant struct {
	Name        string `json:"name"`
	Founder     string `json:"founder"`
	Location    string `json:"location"`
	Country     string `json:"country"`
	Fundation   string `json:"fundation"`
	Headquarter string `json:"headquarter"`
}

func UnmarshalRestaurant(msg []byte) *Restaurant {
	var restaurant Restaurant

	if err := json.Unmarshal(msg, &restaurant); err != nil {
		zap.Log.Error("Failed to make unmarshal to restaurant: ", err)
		return nil
	}

	return &restaurant
}
