package branch

import (
	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/models"
	"github.com/gofiber/fiber/v2/log"
)

type BranchReportPDF struct {
}

func NewBranchReport() IReport {
	return &BranchReportPDF{}
}

func (br *BranchReportPDF) GenerateReport(msg *consumer.Message) {
	branches := models.UnmarshalBranches(msg.Data)

	for _, b := range *branches {
		log.Info(b)
	}
}
