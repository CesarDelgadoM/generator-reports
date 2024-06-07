package generators

import "github.com/CesarDelgadoM/generator-reports/internal/consumer"

const (
	// Formats
	PDF   = "PDF"
	EXCEL = "EXCEL"
)

// Interface to strategies
type IReport interface {
	GenerateReport(msg *consumer.Message)
}
