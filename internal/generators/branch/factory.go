package branch

import (
	"errors"

	"github.com/CesarDelgadoM/generator-reports/config"
	"github.com/CesarDelgadoM/generator-reports/internal/generators"
	"github.com/CesarDelgadoM/generator-reports/pkg/logger/zap"
)

func GetGeneratorBranchReport(format string, config *config.Config) (generators.IReport, error) {

	switch format {
	case generators.PDF:
		zap.Log.Info("Pdf format")
		return NewBranchReportPdf(config.Branch.Pdf), nil

	case generators.EXCEL:
		zap.Log.Info("Excel format")
		return nil, nil

	default:
		return nil, errors.New("File format not implemented: " + format)
	}
}
