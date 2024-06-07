package branch

import (
	"github.com/CesarDelgadoM/generator-reports/internal/consumer"
	"github.com/CesarDelgadoM/generator-reports/internal/generators"
	"github.com/CesarDelgadoM/generator-reports/internal/models"
	"github.com/gofiber/fiber/v2/log"
)

type BranchReportPdf struct {
}

func NewBranchReportPdf() generators.IReport {
	return &BranchReportPdf{}
}

func (br *BranchReportPdf) GenerateReport(msg *consumer.Message) {
	branches := models.UnmarshalBranches(msg.Data)

	for _, b := range *branches {
		log.Info(b)
	}
}
