package restaurant

import (
	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/generators"
	"github.com/CesarDelgadoM/generator-reports/internal/models"
	"github.com/gofiber/fiber/v2/log"
)

type RestaurantReportPDF struct {
}

func NewRestaurantReport() generators.IReport {
	return &RestaurantReportPDF{}
}

func (rr *RestaurantReportPDF) GenerateReport(msg *consumer.Message) {
	restaurant := models.UnmarshalRestaurant(msg.Data)

	log.Info(restaurant)
}
